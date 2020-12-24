package app

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	
	"bean/components/scalar"
	"bean/components/util"
	"bean/components/util/connect"
	"bean/pkg/app/model/dto"
)

func bundle() *AppBundle {
	idr := util.MockIdentifier()
	log := util.MockLogger()
	bun, _ := NewApplicationBundle(idr, log, nil, nil)
	
	return bun
}

func Test(t *testing.T) {
	ass := assert.New(t)
	bun := bundle()
	con := util.MockDatabase().WithContext(context.Background())
	ctx := connect.DBToContext(context.Background(), con)
	util.MockInstall(bun, con)
	
	t.Run("database schema", func(t *testing.T) {
		ass.True(
			con.Migrator().HasTable("applications"),
			"Table `applications` was created",
		)
	})
	
	// create
	oCreate, err := bun.Service.Create(ctx, &dto.ApplicationCreateInput{IsActive: false, Title: scalar.NilString("QA app")})
	ass.NoError(err)
	ass.NotNil(oCreate)
	ass.Equal(false, oCreate.App.IsActive)
	
	t.Run("update", func(t *testing.T) {
		t.Run("useless input", func(t *testing.T) {
			oUpdate, err := bun.Service.Update(ctx, &dto.ApplicationUpdateInput{
				Id:      oCreate.App.ID,
				Version: oCreate.App.Version,
			})
			
			ass.Error(err)
			ass.Nil(oUpdate)
		})
		
		t.Run("status", func(t *testing.T) {
			app, _ := bun.Service.Load(ctx, oCreate.App.ID)
			oUpdate, err := bun.Service.Update(ctx, &dto.ApplicationUpdateInput{
				Id:       app.ID,
				Version:  app.Version,
				IsActive: scalar.NilBool(true),
			})
			
			ass.NoError(err)
			ass.NotNil(oUpdate)
			ass.Equal(true, oUpdate.App.IsActive)
		})
	})
	
	t.Run("delete", func(t *testing.T) {
		app, _ := bun.Service.Load(ctx, oCreate.App.ID)
		now := time.Now()
		oDelete, err := bun.Service.Delete(ctx, dto.ApplicationDeleteInput{
			Id:      app.ID,
			Version: app.Version,
		})
		
		ass.NoError(err)
		ass.NotNil(oDelete)
		ass.True(now.UnixNano() <= oDelete.App.DeletedAt.UnixNano())
	})
}
