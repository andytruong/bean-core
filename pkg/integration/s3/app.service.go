package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/components/util"
	"bean/components/util/connect"
	model2 "bean/pkg/app/model"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type ApplicationService struct {
	bundle *S3Bundle
}

func (service *ApplicationService) Load(ctx context.Context, id string) (*model2.Application, error) {
	app := &model2.Application{}

	// TODO: don't allow to load pending deleted S3 applications.
	err := service.bundle.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&app).
		Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (service *ApplicationService) Create(ctx context.Context, in *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	var app *model2.Application

	err := connect.Transaction(
		ctx,
		service.bundle.db,
		func(tx *gorm.DB) error {
			app = &model2.Application{
				ID:        service.bundle.id.MustULID(),
				Version:   service.bundle.id.MustULID(),
				IsActive:  in.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			}

			err := tx.Create(&app).Error
			if nil != err {
				return err
			} else if err := service.bundle.credentialService.onAppCreate(tx, app, in.Credentials); nil != err {
				return err
			} else if err = service.bundle.policyService.onAppCreate(tx, app, in.Policies); nil != err {
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

func (service *ApplicationService) Update(ctx context.Context, in *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
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

	if deletedAt, ok := ctx.Value(model.DeleteContextKey).(time.Time); ok {
		app.DeletedAt = &deletedAt
		changed = true
	}

	if !changed {
		if nil == in.Credentials && nil == in.Policies {
			return nil, util.ErrorUselessInput
		}
	}

	app.Version = service.bundle.id.MustULID()
	app.UpdatedAt = time.Now()
	err = connect.Transaction(
		ctx,
		service.bundle.db,
		func(tx *gorm.DB) error {
			err := tx.Save(&app).Error
			if nil != err {
				return err
			}

			err = service.bundle.credentialService.onAppUpdate(tx, app, in.Credentials)
			if nil != err {
				return err
			}

			err = service.bundle.policyService.onAppUpdate(tx, app, in.Policies)
			if nil != err {
				return err
			}

			return nil
		},
	)

	return &dto.S3ApplicationMutationOutcome{App: app, Errors: nil}, err
}

func (service *ApplicationService) Delete(ctx context.Context, in dto.S3ApplicationDeleteInput) (*dto.S3ApplicationMutationOutcome, error) {
	ctx = context.WithValue(ctx, model.DeleteContextKey, time.Now())

	return service.Update(ctx, &dto.S3ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: scalar.NilBool(true),
	})
}

func (service *ApplicationService) S3UploadToken(ctx context.Context, in dto.S3UploadTokenInput) (map[string]interface{}, error) {
	// get claims from context
	claims, ok := ctx.Value(claim.ClaimsContextKey).(*claim.Payload)
	if !ok {
		return nil, util.ErrorAuthRequired
	}

	// load application
	app, err := service.bundle.AppService.Load(ctx, in.ApplicationId)
	if nil != err {
		return nil, err
	} else {
		cred, err := service.bundle.credentialService.loadByApplicationId(ctx, in.ApplicationId)
		if nil != err {
			return nil, err
		} else if client, err := service.bundle.credentialService.client(cred); nil != err {
			return nil, err
		} else {
			policy := minio.NewPostPolicy()

			err := scalar.NoError(
				policy.SetBucket(cred.Bucket),
				policy.SetKey(string(in.FilePath)),
				policy.SetExpires(time.Now().UTC().Add(4*time.Hour)),
				policy.SetContentType(string(in.ContentType)),
				policy.SetUserMetadata("app", app.ID),
				policy.SetUserMetadata("sid", claims.SessionId()),
				policy.SetUserMetadata("nid", claims.SpaceId()),
				policy.SetContentLengthRange(1, 10*1024*1024), // TODO: generate per application's policy
			)

			if nil != err {
				return nil, err
			} else if _, formData, err := client.PresignedPostPolicy(ctx, policy); nil != err {
				return nil, err
			} else {
				response := map[string]interface{}{}

				for k, v := range formData {
					response[k] = v
				}

				return response, nil
			}
		}
	}
}
