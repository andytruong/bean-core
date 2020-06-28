package s3

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
)

func bean() *S3IntegrationBean {
	db := util.MockDatabase()
	id := util.MockIdentifier()
	logger := util.MockLogger()
	bean := NewS3Integration(db, id, logger, &Genetic{Key: []byte("01EBWB516AP6BQD7")})
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
		t.Run("Crentials", func(t *testing.T) {
			t.Run("Encrypt", func(t *testing.T) {
				encrypted := this.coreCredentials.encrypt("xxxxxxxxxxxxxxxxxxxxx")
				decrypted := this.coreCredentials.decrypt(encrypted)

				ass.Equal("xxxxxxxxxxxxxxxxxxxxx", decrypted)
				ass.True(len(encrypted)*2 <= 256)
			})
		})

		t.Run("CRUD", func(t *testing.T) {
			oCreate, err := this.coreApp.Create(ctx, dto.S3ApplicationCreateInput{
				Slug:     "qa",
				IsActive: false,
				Credentials: dto.S3ApplicationCredentialsCreateInput{
					Endpoint:  "http://localhost:9090",
					IsSecure:  false,
					AccessKey: "minio",
					SecretKey: "minio",
				},
			})

			ass.NoError(err)
			ass.NotNil(oCreate)
			ass.Equal(false, oCreate.App.IsActive)
			ass.Equal("qa", oCreate.App.Slug)

			t.Run("Update", func(t *testing.T) {
				t.Run("Useless input", func(t *testing.T) {
					oUpdate, err := this.coreApp.Update(ctx, dto.S3ApplicationUpdateInput{
						Id:      oCreate.App.ID,
						Version: oCreate.App.Version,
					})

					ass.Error(err)
					ass.Nil(oUpdate)
				})

				t.Run("Status", func(t *testing.T) {
					app, _ := this.coreApp.Load(ctx, oCreate.App.ID)
					oUpdate, err := this.coreApp.Update(ctx, dto.S3ApplicationUpdateInput{
						Id:       app.ID,
						Version:  app.Version,
						IsActive: util.NilBool(true),
					})

					ass.NoError(err)
					ass.NotNil(oUpdate)
					ass.Equal(true, oUpdate.App.IsActive)
				})

				t.Run("Slug", func(t *testing.T) {
					app, _ := this.coreApp.Load(ctx, oCreate.App.ID)
					oUpdate, err := this.coreApp.Update(ctx, dto.S3ApplicationUpdateInput{
						Id:      app.ID,
						Version: app.Version,
						Slug:    util.NilString("test"),
					})

					ass.NoError(err)
					ass.NotNil(oUpdate)
					ass.Equal("test", oUpdate.App.Slug)
				})

				t.Run("Credentials", func(t *testing.T) {
					app, _ := this.coreApp.Load(ctx, oCreate.App.ID)
					oUpdate, err := this.coreApp.Update(ctx, dto.S3ApplicationUpdateInput{
						Id:      app.ID,
						Version: app.Version,
						Credentials: &dto.S3ApplicationCredentialsUpdateInput{
							Endpoint:  util.NilUri("http://localhost:9191"),
							IsSecure:  util.NilBool(false),
							AccessKey: util.NilString("minio"),
							SecretKey: util.NilString("minio"),
						},
					})

					ass.NoError(err)
					ass.NotNil(oUpdate)
				})
			})

			t.Run("delete", func(t *testing.T) {
				app, _ := this.coreApp.Load(ctx, oCreate.App.ID)
				now := time.Now()
				oDelete, err := this.coreApp.Delete(ctx, dto.S3ApplicationDeleteInput{
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
