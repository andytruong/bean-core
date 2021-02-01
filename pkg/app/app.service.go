package app

import (
	"context"
	"time"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app/model"
	"bean/pkg/app/model/dto"
)

type AppService struct {
	bundle *Bundle
}

func (srv *AppService) Load(ctx context.Context, id string) (*model.Application, error) {
	con := connect.DB(ctx)
	app := &model.Application{}
	err := con.Where("id = ? AND deleted_at IS NULL", id).Take(&app).Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (srv *AppService) Create(ctx context.Context, in *dto.ApplicationCreateInput) (*dto.ApplicationOutcome, error) {
	claims := claim.ContextToPayload(ctx)
	if nil == claims {
		return nil, util.ErrorAuthRequired
	}

	app := &model.Application{
		ID:        srv.bundle.idr.ULID(),
		Version:   srv.bundle.idr.ULID(),
		SpaceId:   claims.SpaceId(),
		IsActive:  in.IsActive,
		Title:     in.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	err := connect.DB(ctx).Create(&app).Error
	if nil != err {
		return nil, err
	}

	return &dto.ApplicationOutcome{App: app, Errors: nil}, nil
}

func (srv *AppService) Update(ctx context.Context, in *dto.ApplicationUpdateInput) (*dto.ApplicationOutcome, error) {
	con := connect.DB(ctx)
	app, err := srv.Load(ctx, in.Id)

	if nil != err {
		return nil, err
	} else if app.Version != in.Version {
		return nil, util.ErrorVersionConflict
	} else if app.DeletedAt != nil {
		return nil, util.ErrorLocked
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

	app.Version = srv.bundle.idr.ULID()
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
