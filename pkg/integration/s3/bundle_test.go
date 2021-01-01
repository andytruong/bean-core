package s3

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app"
	"bean/pkg/config"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

func bundle() *S3Bundle {
	idr := util.MockIdentifier()
	log := util.MockLogger()
	appBundle, _ := app.NewApplicationBundle(idr, log, nil, nil)
	configBundle := config.NewConfigBundle(idr, log)
	bun := NewS3Integration(idr, log, &S3Configuration{Key: "01EBWB516AP6BQD7"}, appBundle, configBundle)

	return bun
}

func Test(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	t.Run("DB schema", func(t *testing.T) {
		ass.True(db.Migrator().HasTable("s3_application_policy"))
	})

	t.Run("Service", func(t *testing.T) {
		t.Run("Credentials", func(t *testing.T) {
			t.Run("Encrypt", func(t *testing.T) {
				encrypted := bundle.credentialService.encrypt("xxxxxxxxxxxxxxxxxxxxx")
				decrypted := bundle.credentialService.decrypt(encrypted)

				ass.Equal("xxxxxxxxxxxxxxxxxxxxx", decrypted)
				ass.True(len(encrypted)*2 <= 256)
			})
		})

		t.Run("CRUD", func(t *testing.T) {
			oCreate, err := bundle.AppService.Create(ctx, &dto.S3ApplicationCreateInput{
				IsActive: false,
				Credentials: dto.S3ApplicationCredentialsCreateInput{
					Endpoint:  "http://localhost:9000",
					Bucket:    "test",
					IsSecure:  false,
					AccessKey: "minioadmin",
					SecretKey: "minioadmin",
				},
				Policies: []dto.S3ApplicationPolicyCreateInput{
					{
						Kind:  model.PolicyKindFileExtensions,
						Value: "jpeg gif png webp",
					},
					{
						Kind:  model.PolicyKindRateLimit,
						Value: "1MB/user/hour",
					},
					{
						Kind:  model.PolicyKindRateLimit,
						Value: "1GB/space/hour",
					},
				},
			})

			ass.NoError(err)
			ass.NotNil(oCreate)
			ass.Equal(false, oCreate.App.IsActive)

			t.Run("policies", func(t *testing.T) {
				policies := []model.Policy{}
				err := db.Where("application_id = ?", oCreate.App.ID).Find(&policies).Error
				ass.NoError(err)
				ass.Equal(3, len(policies))
				ass.Equal(policies[0].Kind, model.PolicyKindFileExtensions)
				ass.Equal(policies[1].Kind, model.PolicyKindRateLimit)
				ass.Equal(policies[2].Kind, model.PolicyKindRateLimit)
				ass.Equal(policies[0].Value, "jpeg gif png webp")
				ass.Equal(policies[1].Value, "1MB/user/hour")
				ass.Equal(policies[2].Value, "1GB/space/hour")
			})

			t.Run("update", func(t *testing.T) {
				t.Run("credentials", func(t *testing.T) {
					app, _ := bundle.appBundle.Service.Load(ctx, oCreate.App.ID)
					oUpdate, err := bundle.AppService.Update(ctx, &dto.S3ApplicationUpdateInput{
						Id:      app.ID,
						Version: app.Version,
						Credentials: &dto.S3ApplicationCredentialsUpdateInput{
							Endpoint:  scalar.NilUri("http://localhost:9191"),
							Bucket:    scalar.NilString("test"),
							IsSecure:  scalar.NilBool(false),
							AccessKey: scalar.NilString("minio"),
							SecretKey: scalar.NilString("minio"),
						},
					})

					ass.NoError(err)
					ass.NotNil(oUpdate)

					// reload & assert
					{
						cred, err := bundle.credentialService.loadByApplicationId(ctx, app.ID)
						ass.NoError(err)
						ass.Equal("http://localhost:9191", string(cred.Endpoint))
						ass.Equal("test", cred.Bucket)
						ass.Equal("minio", cred.AccessKey)
						ass.NotEqual("minio", cred.SecretKey, "value is encrypted")
						ass.Equal(false, cred.IsSecure)
					}
				})

				t.Run("policies", func(t *testing.T) {
					app, _ := bundle.appBundle.Service.Load(ctx, oCreate.App.ID)
					policies := []model.Policy{}
					err := db.Where("application_id = ?", oCreate.App.ID).Find(&policies).Error
					ass.NoError(err)

					// before update
					{
						ass.Equal(3, len(policies))
						ass.Equal(policies[0].Kind, model.PolicyKindFileExtensions)
						ass.Equal(policies[1].Kind, model.PolicyKindRateLimit)
						ass.Equal(policies[2].Kind, model.PolicyKindRateLimit)
						ass.Equal(policies[0].Value, "jpeg gif png webp")
						ass.Equal(policies[1].Value, "1MB/user/hour")
						ass.Equal(policies[2].Value, "1GB/space/hour")
					}

					oUpdate, err := bundle.AppService.Update(ctx, &dto.S3ApplicationUpdateInput{
						Id:      app.ID,
						Version: app.Version,
						Policies: &dto.S3ApplicationPolicyMutationInput{
							Create: []dto.S3ApplicationPolicyCreateInput{
								{
									Kind:  model.PolicyKindFileExtensions,
									Value: "raw",
								},
							},
							Update: []dto.S3ApplicationPolicyUpdateInput{
								{
									Id:    policies[1].ID,
									Value: "2MB/user/hour",
								},
							},
							Delete: []dto.S3ApplicationPolicyDeleteInput{
								{
									Id: policies[2].ID,
								},
							},
						},
					})

					ass.NoError(err)
					ass.NotNil(oUpdate)

					// after update: add 1, update 1, remove 1
					{
						policies, err := bundle.policyService.loadByApplicationId(ctx, app.ID)
						ass.NoError(err)
						ass.Equal(3, len(policies))
						ass.Equal(policies[0].Kind, model.PolicyKindFileExtensions)
						ass.Equal(policies[1].Kind, model.PolicyKindRateLimit)
						ass.Equal(policies[2].Kind, model.PolicyKindFileExtensions)
						ass.Equal(policies[0].Value, "jpeg gif png webp")
						ass.Equal(policies[1].Value, "2MB/user/hour")
						ass.Equal(policies[2].Value, "raw")
					}
				})
			})
		})
	})
}

func Test_UploadToken(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	bundle.credentialService.transport = connect.MockRoundTrip{
		Callback: func(request *http.Request) (*http.Response, error) {
			response := &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
			}

			content := `<?xml version="1.0" encoding="UTF-8"?>`
			content += `<LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">Europe</LocationConstraint>`
			body := strings.NewReader(content)
			response.Body = ioutil.NopCloser(body)

			return response, nil
		},
	}
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	ctx = context.WithValue(ctx, claim.ClaimsContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{
			Audience: bundle.idr.MustULID(),
			Id:       bundle.idr.MustULID(),
			Subject:  bundle.idr.MustULID(),
		},
		Kind: claim.KindAuthenticated,
	})

	oCreate, err := bundle.AppService.Create(ctx, &dto.S3ApplicationCreateInput{
		IsActive: false,
		Credentials: dto.S3ApplicationCredentialsCreateInput{
			Endpoint:  "http://localhost:9000",
			Bucket:    "test",
			IsSecure:  false,
			AccessKey: "minioadmin",
			SecretKey: "minioadmin",
		},
		Policies: []dto.S3ApplicationPolicyCreateInput{
			{
				Kind:  model.PolicyKindFileExtensions,
				Value: "jpeg gif png webp",
			},
			{
				Kind:  model.PolicyKindRateLimit,
				Value: "1MB/user/hour",
			},
			{
				Kind:  model.PolicyKindRateLimit,
				Value: "1GB/space/hour",
			},
		},
	})

	ass.NoError(err)
	ass.NotNil(oCreate)

	formData, err := bundle.AppService.S3UploadToken(ctx, dto.S3UploadTokenInput{
		ApplicationId: oCreate.App.ID,
		FilePath:      "/path/to/image.png",
		ContentType:   scalar.ImagePNG,
	})

	ass.NoError(err)
	ass.Equal(formData["bucket"], "test")
	ass.Equal(formData["key"], "/path/to/image.png")
	ass.Equal(formData["Content-Type"], string(scalar.ImagePNG))
	ass.NotEmpty(formData["policy"])
	ass.NotEmpty(formData["x-amz-credential"])
	ass.NotEmpty(formData["x-amz-meta-nid"])
	ass.NotEmpty(formData["x-amz-signature"])
}
