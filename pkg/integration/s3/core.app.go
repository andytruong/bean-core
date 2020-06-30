package s3

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type CoreApplication struct {
	bean     *S3IntegrationBean
	Resolver *ApplicationResolver
}

func (this *CoreApplication) Load(ctx context.Context, id string) (*model.Application, error) {
	app := &model.Application{}

	// TODO: don't allow to load pending deleted S3 applications.
	err := this.bean.db.
		WithContext(ctx).
		Table(connect.TableIntegrationS3App).
		Where("id = ?", id).
		First(&app).
		Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (this *CoreApplication) Create(ctx context.Context, in *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	var app *model.Application

	err := connect.Transaction(
		ctx,
		this.bean.db,
		func(tx *gorm.DB) error {
			app = &model.Application{
				Slug:      in.Slug,
				ID:        this.bean.id.MustULID(),
				Version:   this.bean.id.MustULID(),
				IsActive:  in.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			}

			err := tx.Table(connect.TableIntegrationS3App).Create(&app).Error
			if nil != err {
				return err
			} else if err := this.bean.coreCredentials.onAppCreate(tx, app, in.Credentials); nil != err {
				return err
			} else if err = this.bean.corePolicy.onAppCreate(tx, app, in.Policies); nil != err {
				return err
			}

			return nil
		},
	)

	if nil != err {
		return nil, err
	}

	return &dto.S3ApplicationMutationOutcome{App: app, Errors: nil}, nil
}

func (this *CoreApplication) Update(ctx context.Context, in dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
	app, err := this.Load(ctx, in.Id)

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

	if nil != in.Slug {
		if app.Slug != *in.Slug {
			app.Slug = *in.Slug
			changed = true
		}
	}

	if deletedAt, ok := ctx.Value("bean.integration-s3.delete").(time.Time); ok {
		app.DeletedAt = &deletedAt
		changed = true
	}

	if !changed {
		if nil == in.Credentials && nil == in.Polices {
			return nil, util.ErrorUselessInput
		}
	}

	app.Version = this.bean.id.MustULID()
	app.UpdatedAt = time.Now()
	err = connect.Transaction(
		ctx,
		this.bean.db,
		func(tx *gorm.DB) error {
			err := tx.Table(connect.TableIntegrationS3App).Save(&app).Error
			if nil != err {
				return err
			}

			err = this.bean.coreCredentials.onAppUpdate(tx, app, in.Credentials)
			if nil != err {
				return err
			}

			err = this.bean.corePolicy.onAppUpdate(tx, app, in.Polices)
			if nil != err {
				return err
			}

			return nil
		},
	)

	return &dto.S3ApplicationMutationOutcome{App: app, Errors: nil}, err
}

func (this *CoreApplication) Delete(ctx context.Context, in dto.S3ApplicationDeleteInput) (*dto.S3ApplicationMutationOutcome, error) {
	ctx = context.WithValue(ctx, "bean.integration-s3.delete", time.Now())

	return this.Update(ctx, dto.S3ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: util.NilBool(true),
	})
}
