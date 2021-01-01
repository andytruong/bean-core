package app

import (
	"context"
	"time"

	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app/model"
	"bean/pkg/app/model/dto"
)

type AppService struct {
	bundle *AppBundle
}

func (srv *AppService) Load(ctx context.Context, id string) (*model.Application, error) {
	con := connect.ContextToDB(ctx)
	app := &model.Application{}

	// TODO: don't allow to load pending deleted applications.
	err := con.Where("id = ?", id).First(&app).Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (srv *AppService) Create(ctx context.Context, in *dto.ApplicationCreateInput) (*dto.ApplicationOutcome, error) {
	app := &model.Application{
		ID:        srv.bundle.idr.MustULID(),
		Version:   srv.bundle.idr.MustULID(),
		IsActive:  in.IsActive,
		Title:     in.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	err := connect.ContextToDB(ctx).Create(&app).Error
	if nil != err {
		return nil, err
	}

	return &dto.ApplicationOutcome{App: app, Errors: nil}, nil
}

func (srv *AppService) Update(ctx context.Context, in *dto.ApplicationUpdateInput) (*dto.ApplicationOutcome, error) {
	con := connect.ContextToDB(ctx)
	app, err := srv.Load(ctx, in.Id)

	if nil != err {
		return nil, err
	} else if app.Version != in.Version {
		return nil, util.ErrorVersionConflict
	}

	changed := false
	if nil != in.IsActive {
		if app.IsActive != *in.IsActive {
			app.IsActive = *in.IsActive
			changed = true
		}
	}

	if nil != in.Title {
		if *app.Title != *in.Title {
			app.Title = in.Title
			changed = true
		}
	}

	if deletedAt, ok := ctx.Value(model.DeleteContextKey).(time.Time); ok {
		app.DeletedAt = &deletedAt
		changed = true
	}

	if !changed {
		return &dto.ApplicationOutcome{App: app, Errors: nil}, util.ErrorUselessInput
	}

	app.Version = srv.bundle.idr.MustULID()
	app.UpdatedAt = time.Now()
	err = con.Save(&app).Error

	return &dto.ApplicationOutcome{App: app, Errors: nil}, err
}

func (srv *AppService) Delete(ctx context.Context, in dto.ApplicationDeleteInput) (*dto.ApplicationOutcome, error) {
	ctx = context.WithValue(ctx, model.DeleteContextKey, time.Now())

	return srv.Update(ctx, &dto.ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: scalar.NilBool(true),
	})
}
