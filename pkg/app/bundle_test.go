package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app/model/dto"
)

func bundle() *Bundle {
	idr := util.MockIdentifier()
	log := util.MockLogger()
	hook := module.NewHook()
	bun, _ := NewApplicationBundle(idr, log, hook, nil, nil)

	return bun
}

func Test(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	claims := claim.NewPayload()
	claims.SetSpaceId(bundle.idr.ULID())
	ctx = claim.PayloadToContext(ctx, &claims)
	connect.MockInstall(ctx, bundle)

	t.Run("database schema", func(t *testing.T) {
		ass.True(
			db.Migrator().HasTable("applications"),
			"Table `applications` was created",
		)
	})

	// create
	oCreate, err := bundle.Service.Create(ctx, &dto.ApplicationCreateInput{
		IsActive: false, Title: scalar.NilString("QA app"),
	})
	ass.NoError(err)
	ass.NotNil(oCreate)
	ass.Equal(false, oCreate.App.IsActive)

	t.Run("update", func(t *testing.T) {
		t.Run("useless input", func(t *testing.T) {
			oUpdate, err := bundle.Service.Update(ctx, &dto.ApplicationUpdateInput{
				Id:      oCreate.App.ID,
				Version: oCreate.App.Version,
				Title:   scalar.NilString("QA app"),
			})

			ass.Error(err)
			ass.Equal(oUpdate.App.Version, oCreate.App.Version)
		})

		t.Run("status", func(t *testing.T) {
			app, _ := bundle.Service.Load(ctx, oCreate.App.ID)
			oUpdate, err := bundle.Service.Update(ctx, &dto.ApplicationUpdateInput{
				Id:       app.ID,
				Version:  app.Version,
				IsActive: scalar.NilBool(true),
			})

			ass.NoError(err)
			ass.NotNil(oUpdate)
			ass.NotEqual(oUpdate.App.Version, app.Version)
			ass.Equal(true, oUpdate.App.IsActive)
		})

		t.Run("title", func(t *testing.T) {
			app, _ := bundle.Service.Load(ctx, oCreate.App.ID)
			oUpdate, err := bundle.Service.Update(ctx, &dto.ApplicationUpdateInput{
				Id:      app.ID,
				Version: app.Version,
				Title:   scalar.NilString("Quality Assurance application"),
			})

			ass.NoError(err)
			ass.NotNil(oUpdate)
			ass.NotEqual(oUpdate.App.Version, app.Version)
			ass.Equal(*oUpdate.App.Title, "Quality Assurance application")
		})
	})

	t.Run("delete", func(t *testing.T) {
		app, _ := bundle.Service.Load(ctx, oCreate.App.ID)
		now := time.Now()

		oDelete, err := bundle.Service.Delete(ctx, dto.ApplicationDeleteInput{
			Id:      app.ID,
			Version: app.Version,
		})

		ass.NoError(err)
		ass.NotNil(oDelete)
		ass.True(now.UnixNano() <= oDelete.App.DeletedAt.UnixNano())

		t.Run("deleted application", func(t *testing.T) {
			oDelete, err := bundle.Service.Delete(ctx, dto.ApplicationDeleteInput{
				Id:      app.ID,
				Version: oDelete.App.Version,
			})

			ass.Error(err)
			ass.Equal(gorm.ErrRecordNotFound, err)
			ass.Nil(oDelete)
		})
	})
}
