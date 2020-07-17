package s3

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

func bean() *S3IntegrationBean {
	db := util.MockDatabase()
	id := util.MockIdentifier()
	logger := util.MockLogger()
	bean := NewS3Integration(db, id, logger, &Genetic{Key: "01EBWB516AP6BQD7"})
	util.MockInstall(bean, db)

	return bean
}

func Test(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	ctx := context.Background()

	t.Run("DB schema", func(t *testing.T) {
		this.db.Migrator().HasTable("s3_application")
	})

	t.Run("Core Application", func(t *testing.T) {
		t.Run("Credentials", func(t *testing.T) {
			t.Run("Encrypt", func(t *testing.T) {
				encrypted := this.coreCredentials.encrypt("xxxxxxxxxxxxxxxxxxxxx")
				decrypted := this.coreCredentials.decrypt(encrypted)

				ass.Equal("xxxxxxxxxxxxxxxxxxxxx", decrypted)
				ass.True(len(encrypted)*2 <= 256)
			})
		})

		t.Run("CRUD", func(t *testing.T) {
			oCreate, err := this.CoreApp.Create(ctx, &dto.S3ApplicationCreateInput{
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
						Value: "1GB/namespace/hour",
					},
				},
			})

			ass.NoError(err)
			ass.NotNil(oCreate)
			ass.Equal(false, oCreate.App.IsActive)

			t.Run("policies", func(t *testing.T) {
				policies := []model.Policy{}
				err := this.db.
					Table(connect.TableIntegrationS3Policy).
					Where("application_id = ?", oCreate.App.ID).
					Find(&policies).
					Error
				ass.NoError(err)
				ass.Equal(3, len(policies))
				ass.Equal(policies[0].Kind, model.PolicyKindFileExtensions)
				ass.Equal(policies[1].Kind, model.PolicyKindRateLimit)
				ass.Equal(policies[2].Kind, model.PolicyKindRateLimit)
				ass.Equal(policies[0].Value, "jpeg gif png webp")
				ass.Equal(policies[1].Value, "1MB/user/hour")
				ass.Equal(policies[2].Value, "1GB/namespace/hour")
			})

			t.Run("Update", func(t *testing.T) {
				t.Run("Useless input", func(t *testing.T) {
					oUpdate, err := this.CoreApp.Update(ctx, &dto.S3ApplicationUpdateInput{
						Id:      oCreate.App.ID,
						Version: oCreate.App.Version,
					})

					ass.Error(err)
					ass.Nil(oUpdate)
				})

				t.Run("Status", func(t *testing.T) {
					app, _ := this.CoreApp.Load(ctx, oCreate.App.ID)
					oUpdate, err := this.CoreApp.Update(ctx, &dto.S3ApplicationUpdateInput{
						Id:       app.ID,
						Version:  app.Version,
						IsActive: scalar.NilBool(true),
					})

					ass.NoError(err)
					ass.NotNil(oUpdate)
					ass.Equal(true, oUpdate.App.IsActive)
				})

				t.Run("Bucket", func(t *testing.T) {
					app, _ := this.CoreApp.Load(ctx, oCreate.App.ID)
					oUpdate, err := this.CoreApp.Update(ctx, &dto.S3ApplicationUpdateInput{
						Id:      app.ID,
						Version: app.Version,
					})

					ass.Error(err)
					ass.Equal(err, util.ErrorUselessInput)
					ass.Nil(oUpdate)
				})

				t.Run("Credentials", func(t *testing.T) {
					app, _ := this.CoreApp.Load(ctx, oCreate.App.ID)
					oUpdate, err := this.CoreApp.Update(ctx, &dto.S3ApplicationUpdateInput{
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
						cred, err := this.CoreApp.Resolver.Credentials(ctx, app)
						ass.NoError(err)
						ass.Equal("http://localhost:9191", string(cred.Endpoint))
						ass.Equal("test", cred.Bucket)
						ass.Equal("minio", cred.AccessKey)
						ass.NotEqual("minio", cred.SecretKey, "value is encrypted")
						ass.Equal(false, cred.IsSecure)
					}
				})

				t.Run("Policies", func(t *testing.T) {
					app, _ := this.CoreApp.Load(ctx, oCreate.App.ID)
					policies := []model.Policy{}
					err := this.db.Table(connect.TableIntegrationS3Policy).Where("application_id = ?", oCreate.App.ID).Find(&policies).Error
					ass.NoError(err)

					// before update
					{
						ass.Equal(3, len(policies))
						ass.Equal(policies[0].Kind, model.PolicyKindFileExtensions)
						ass.Equal(policies[1].Kind, model.PolicyKindRateLimit)
						ass.Equal(policies[2].Kind, model.PolicyKindRateLimit)
						ass.Equal(policies[0].Value, "jpeg gif png webp")
						ass.Equal(policies[1].Value, "1MB/user/hour")
						ass.Equal(policies[2].Value, "1GB/namespace/hour")
					}

					oUpdate, err := this.CoreApp.Update(ctx, &dto.S3ApplicationUpdateInput{
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
						policies, err := this.CoreApp.Resolver.Polices(ctx, app)
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

			t.Run("delete", func(t *testing.T) {
				app, _ := this.CoreApp.Load(ctx, oCreate.App.ID)
				now := time.Now()
				oDelete, err := this.CoreApp.Delete(ctx, dto.S3ApplicationDeleteInput{
					Id:      app.ID,
					Version: app.Version,
				})

				ass.NoError(err)
				ass.NotNil(oDelete)
				ass.True(now.UnixNano() <= oDelete.App.DeletedAt.UnixNano())
			})
		})
	})
}

func Test_UploadToken(t *testing.T) {
	ass := assert.New(t)
	this := bean()

	this.coreCredentials.transport = util.MockRoundTrip{
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

	ctx := context.WithValue(context.Background(), claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{
			Audience: this.id.MustULID(),
			Id:       this.id.MustULID(),
			Subject:  this.id.MustULID(),
		},
		Kind: claim.KindAuthenticated,
	})

	oCreate, err := this.CoreApp.Create(ctx, &dto.S3ApplicationCreateInput{
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
				Value: "1GB/namespace/hour",
			},
		},
	})

	ass.NoError(err)
	ass.NotNil(oCreate)

	formData, err := this.CoreApp.Resolver.S3UploadToken(ctx, dto.S3UploadTokenInput{
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
