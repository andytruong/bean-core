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
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type ApplicationService struct {
	bundle *S3IntegrationBundle
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

func (this *ApplicationService) S3UploadToken(ctx context.Context, in dto.S3UploadTokenInput) (map[string]interface{}, error) {
	// get claims from context
	claims, ok := ctx.Value(claim.ContextKey).(*claim.Payload)
	if !ok {
		return nil, util.ErrorAuthRequired
	}

	// load application
	app, err := this.bundle.AppService.Load(ctx, in.ApplicationId)
	if nil != err {
		return nil, err
	} else {
		cred, err := this.bundle.credentialService.loadByApplicationId(ctx, in.ApplicationId)
		if nil != err {
			return nil, err
		} else if client, err := this.bundle.credentialService.client(cred); nil != err {
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

	return nil, nil
}
