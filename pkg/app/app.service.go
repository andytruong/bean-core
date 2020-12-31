package app

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app/model"
	"bean/pkg/app/model/dto"
)

type AppService struct {
	bundle *AppBundle
}

func (service *AppService) Load(ctx context.Context, id string) (*model.Application, error) {
	con := connect.ContextToDB(ctx)
	app := &model.Application{}

	// TODO: don't allow to load pending deleted applications.
	err := con.Where("id = ?", id).First(&app).Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (service *AppService) Create(ctx context.Context, in *dto.ApplicationCreateInput) (*dto.ApplicationOutcome, error) {
	var app *model.Application

	tx := connect.ContextToDB(ctx)

	err := connect.Transaction(tx, func(tx *gorm.DB) error {
		app = &model.Application{
			ID:        service.bundle.idr.MustULID(),
			Version:   service.bundle.idr.MustULID(),
			IsActive:  in.IsActive,
			Title:     in.Title,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}

		err := tx.Create(&app).Error
		if nil != err {
			return err
		}

		return nil
	})

	if nil != err {
		return nil, err
	}

	return &dto.ApplicationOutcome{App: app, Errors: nil}, nil
}

func (service *AppService) Update(ctx context.Context, in *dto.ApplicationUpdateInput) (*dto.ApplicationOutcome, error) {
	con := connect.ContextToDB(ctx)
	app, err := service.Load(ctx, in.Id)

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

	app.Version = service.bundle.idr.MustULID()
	app.UpdatedAt = time.Now()
	err = con.Save(&app).Error

	return &dto.ApplicationOutcome{App: app, Errors: nil}, err
}

func (service *AppService) Delete(ctx context.Context, in dto.ApplicationDeleteInput) (*dto.ApplicationOutcome, error) {
	ctx = context.WithValue(ctx, model.DeleteContextKey, time.Now())

	return service.Update(ctx, &dto.ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: scalar.NilBool(true),
	})
}
