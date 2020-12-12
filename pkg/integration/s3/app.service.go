package s3

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/components/scalar"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type ApplicationService struct {
	bundle   *S3IntegrationBundle
	Resolver *ApplicationResolver
}

func (this *ApplicationService) Load(ctx context.Context, id string) (*model.Application, error) {
	app := &model.Application{}

	// TODO: don't allow to load pending deleted S3 applications.
	err := this.bundle.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&app).
		Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (this *ApplicationService) Create(ctx context.Context, in *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	var app *model.Application

	err := connect.Transaction(
		ctx,
		this.bundle.db,
		func(tx *gorm.DB) error {
			app = &model.Application{
				ID:        this.bundle.id.MustULID(),
				Version:   this.bundle.id.MustULID(),
				IsActive:  in.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			}

			err := tx.Create(&app).Error
			if nil != err {
				return err
			} else if err := this.bundle.credentialService.onAppCreate(tx, app, in.Credentials); nil != err {
				return err
			} else if err = this.bundle.policyService.onAppCreate(tx, app, in.Policies); nil != err {
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

func (this *ApplicationService) Update(ctx context.Context, in *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
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

	if deletedAt, ok := ctx.Value("bundle.integration-s3.delete").(time.Time); ok {
		app.DeletedAt = &deletedAt
		changed = true
	}

	if !changed {
		if nil == in.Credentials && nil == in.Policies {
			return nil, util.ErrorUselessInput
		}
	}

	app.Version = this.bundle.id.MustULID()
	app.UpdatedAt = time.Now()
	err = connect.Transaction(
		ctx,
		this.bundle.db,
		func(tx *gorm.DB) error {
			err := tx.Save(&app).Error
			if nil != err {
				return err
			}

			err = this.bundle.credentialService.onAppUpdate(tx, app, in.Credentials)
			if nil != err {
				return err
			}

			err = this.bundle.policyService.onAppUpdate(tx, app, in.Policies)
			if nil != err {
				return err
			}

			return nil
		},
	)

	return &dto.S3ApplicationMutationOutcome{App: app, Errors: nil}, err
}

func (this *ApplicationService) Delete(ctx context.Context, in dto.S3ApplicationDeleteInput) (*dto.S3ApplicationMutationOutcome, error) {
	ctx = context.WithValue(ctx, "bundle.integration-s3.delete", time.Now())

	return this.Update(ctx, &dto.S3ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: scalar.NilBool(true),
	})
}
