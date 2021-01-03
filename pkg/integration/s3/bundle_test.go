package s3

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app"
	appDto "bean/pkg/app/model/dto"
	"bean/pkg/config"
	configDto "bean/pkg/config/model/dto"
	"bean/pkg/infra/api"
	"bean/pkg/integration/s3/model/dto"
)

func bundle() *S3Bundle {
	idr := util.MockIdentifier()
	log := util.MockLogger()
	cnf := &S3Configuration{Key: "01EBWB516AP6BQD7"}
	hook := module.NewHook()
	appBundle, _ := app.NewApplicationBundle(idr, log, hook, nil, nil)
	configBundle := config.NewConfigBundle(idr, log)
	bun := NewS3Integration(idr, log, cnf, appBundle, configBundle)

	return bun
}

func Test_Install(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	t.Run("DB schema", func(t *testing.T) {
		ass.True(db.Migrator().HasTable("s3_upload_token"))
		ass.True(db.Migrator().HasTable("s3_file"))
	})

	t.Run("Config buckets", func(t *testing.T) {
		t.Run("credentials", func(t *testing.T) {
			bucket, err := bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: credentialsConfigSlug})
			ass.NoError(err)
			ass.Equal(bucket.Schema, credentialsConfigSchema)
		})

		t.Run("policy", func(t *testing.T) {
			bucket, err := bundle.configBundle.BucketService.Load(ctx, configDto.BucketKey{Slug: policyConfigSlug})
			ass.NoError(err)
			ass.Equal(bucket.Schema, policyConfigSchema)
		})
	})
}

func Test_Credentials(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()

	claims := claim.NewPayload()
	claims.SetUserId(_id)
	ctx := claim.PayloadToContext(context.Background(), &claims)
	ctx = connect.DBToContext(ctx, db)
	connect.MockInstall(ctx, bundle)

	t.Run("encrypt", func(t *testing.T) {
		encrypted := bundle.credentialService.encrypt("xxxxxxxxxxxxxxxxxxxxx")
		decrypted := bundle.credentialService.decrypt(encrypted)

		ass.Equal("xxxxxxxxxxxxxxxxxxxxx", decrypted)
		ass.True(len(encrypted)*2 <= 256)
	})

	// TODO: If application.inactive?
	t.Run("save", func(t *testing.T) {
		oApp, err := bundle.appBundle.Service.Create(ctx, &appDto.ApplicationCreateInput{IsActive: true})
		ass.NoError(err)
		ass.NotNil(oApp)

		cre, err := bundle.credentialService.save(ctx, dto.S3CredentialsInput{
			Version:       "",
			ApplicationId: oApp.App.ID,
			Endpoint:      "http://localhost:9000",
			Bucket:        "test",
			IsSecure:      false,
			AccessKey:     "minioadmin",
			SecretKey:     "minioadmin",
		})

		t.Run("new", func(t *testing.T) {
			// should see error
			ass.NoError(err)
			ass.NotNil(cre)
			ass.NotEmpty(cre.Version)
			ass.Equal(string(cre.Endpoint), "http://localhost:9000")
			ass.False(cre.IsSecure)
			ass.Equal(cre.AccessKey, "minioadmin")
			ass.Equal(cre.SecretKey, "minioadmin")
		})

		t.Run("load", func(t *testing.T) {
			cre, err := bundle.credentialService.load(ctx, oApp.App.ID)

			ass.NoError(err)
			ass.NotNil(cre)
		})

		t.Run("update.useless-input", func(t *testing.T) {
			_, err := bundle.credentialService.save(ctx, dto.S3CredentialsInput{
				Version:       cre.Version,
				ApplicationId: oApp.App.ID,
				Endpoint:      "http://localhost:9000",
				Bucket:        "test",
				IsSecure:      false,
				AccessKey:     "minioadmin",
				SecretKey:     "minioadmin",
			})

			ass.Error(err)
			ass.Equal(util.ErrorUselessInput, err)
		})

		t.Run("update.version.conflict", func(t *testing.T) {
			_, err := bundle.credentialService.save(ctx, dto.S3CredentialsInput{
				Version:       cre.Id,
				ApplicationId: oApp.App.ID,
				Endpoint:      "http://localhost:9000",
				Bucket:        "test",
				IsSecure:      false,
				AccessKey:     "minioadmin.2",
				SecretKey:     "minioadmin.2",
			})

			ass.Error(err)
			ass.Equal(err, util.ErrorVersionConflict)
		})

		t.Run("update.version.valid", func(t *testing.T) {
			cre, err := bundle.credentialService.save(ctx, dto.S3CredentialsInput{
				Version:       cre.Version,
				ApplicationId: oApp.App.ID,
				Endpoint:      "http://localhost:9000",
				Bucket:        "test",
				IsSecure:      false,
				AccessKey:     "minioadmin.2",
				SecretKey:     "minioadmin.2",
			})

			ass.NoError(err)
			ass.Equal(cre.AccessKey, "minioadmin.2")
			ass.Equal(cre.SecretKey, "minioadmin.2")
		})
	})
}

func Test_Policies(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()

	claims := claim.NewPayload()
	claims.SetUserId(_id)
	ctx := claim.PayloadToContext(context.Background(), &claims)
	ctx = connect.DBToContext(ctx, db)
	connect.MockInstall(ctx, bundle)

	t.Run("save", func(t *testing.T) {
		oApp, err := bundle.appBundle.Service.Create(ctx, &appDto.ApplicationCreateInput{IsActive: true})
		ass.NoError(err)
		ass.NotNil(oApp)

		pol, err := bundle.policyService.save(ctx, dto.UploadPolicyInput{
			Version:        "",
			ApplicationId:  oApp.App.ID,
			FileExtensions: []api.FileType{"jpeg", "gif", "png", "webp"},
			RateLimit: []dto.UploadRateLimitInput{
				{Value: "1MB", Object: "user", Interval: "1 hour"},
				{Value: "1GB", Object: "space", Interval: "1 hour"},
			},
		})

		t.Run("new", func(t *testing.T) {
			ass.NoError(err)
			ass.NotNil(pol)
		})

		t.Run("update.useless-input", func(t *testing.T) {
			_, err := bundle.policyService.save(ctx, dto.UploadPolicyInput{
				Version:        pol.Version,
				ApplicationId:  oApp.App.ID,
				FileExtensions: []api.FileType{"jpeg", "gif", "png", "webp"},
				RateLimit: []dto.UploadRateLimitInput{
					{Value: "1MB", Object: "user", Interval: "1 hour"},
					{Value: "1GB", Object: "space", Interval: "1 hour"},
				},
			})

			ass.Equal(util.ErrorUselessInput, err)
		})

		t.Run("update.version.conflict", func(t *testing.T) {
			_, err := bundle.policyService.save(ctx, dto.UploadPolicyInput{
				Version:        pol.Id,
				ApplicationId:  oApp.App.ID,
				FileExtensions: []api.FileType{"jpeg", "gif", "png", "webp"},
				RateLimit: []dto.UploadRateLimitInput{
					{Value: "2MB", Object: "user", Interval: "1 hour"},
					{Value: "2GB", Object: "space", Interval: "1 hour"},
				},
			})

			ass.Equal(util.ErrorVersionConflict, err)
		})

		t.Run("update.version.valid", func(t *testing.T) {
			newPol, err := bundle.policyService.save(ctx, dto.UploadPolicyInput{
				Version:        pol.Version,
				ApplicationId:  oApp.App.ID,
				FileExtensions: []api.FileType{"jpeg", "gif", "png", "webp"},
				RateLimit: []dto.UploadRateLimitInput{
					{Value: "2MB", Object: "user", Interval: "hour"},
					{Value: "2GB", Object: "space", Interval: "hour"},
				},
			})

			ass.NoError(err)
			ass.NotEqual(pol.Version, newPol.Version)
			ass.Equal("2MB", newPol.RateLimit[0].Value)
			ass.Equal("2GB", newPol.RateLimit[1].Value)
		})
	})
}

func Test_UploadToken(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	bundle.credentialService.transport = connect.MockRoundTrip{
		Callback: func(request *http.Request) (*http.Response, error) {
			response := &http.Response{Status: "OK", StatusCode: http.StatusOK}
			content := `<?xml version="1.0" encoding="UTF-8"?>`
			content += `<LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/">Europe</LocationConstraint>`
			body := strings.NewReader(content)
			response.Body = ioutil.NopCloser(body)

			return response, nil
		},
	}
	db := connect.MockDatabase()
	claims := claim.NewPayload()
	claims.
		SetUserId(_id).
		SetSpaceId(_id).
		SetSessionId(_id)
	ctx := claim.PayloadToContext(context.Background(), &claims)
	ctx = connect.DBToContext(ctx, db)
	connect.MockInstall(ctx, bundle)

	oApp, err := bundle.appBundle.Service.Create(ctx, &appDto.ApplicationCreateInput{IsActive: true})
	ass.NoError(err)
	ass.NotNil(oApp)

	// setup app's settings
	{
		_, err = bundle.credentialService.save(ctx, dto.S3CredentialsInput{
			Version:       "",
			ApplicationId: oApp.App.ID,
			Endpoint:      "http://localhost:9000",
			Bucket:        "test",
			IsSecure:      false,
			AccessKey:     "minioadmin",
			SecretKey:     "minioadmin",
		})
		ass.NoError(err)

		_, err = bundle.policyService.save(ctx, dto.UploadPolicyInput{
			Version:        "",
			ApplicationId:  oApp.App.ID,
			FileExtensions: []api.FileType{"jpeg", "gif", "png", "webp"},
			RateLimit: []dto.UploadRateLimitInput{
				{Value: "1MB", Object: "user", Interval: "1 hour"},
				{Value: "1MB", Object: "space", Interval: "1 hour"},
			},
		})
		ass.NoError(err)
	}

	formData, err := bundle.AppService.CreateUploadToken(ctx, dto.UploadTokenInput{
		ApplicationId: oApp.App.ID,
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
