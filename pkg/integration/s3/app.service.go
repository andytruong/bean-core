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
	dto2 "bean/pkg/app/model/dto"
	"bean/pkg/integration/s3/model/dto"
)

type ApplicationService struct {
	bundle *S3Bundle
}

func (service *ApplicationService) Create(ctx context.Context, in *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	out, err := service.bundle.appBundle.Service.Create(ctx, &dto2.ApplicationCreateInput{IsActive: in.IsActive})
	if nil != err {
		return nil, err
	}

	if nil == out.App {
		return &dto.S3ApplicationMutationOutcome{App: nil, Errors: out.Errors}, nil
	}

	err = connect.Transaction(
		ctx,
		connect.ContextToDB(ctx),
		func(tx *gorm.DB) error {
			ctx := connect.DBToContext(ctx, tx)
			if err := service.bundle.credentialService.onAppCreate(ctx, out.App, in.Credentials); nil != err {
				return err
			} else if err = service.bundle.policyService.onAppCreate(ctx, out.App, in.Policies); nil != err {
				return err
			}

			return nil
		},
	)

	if nil != err {
		return nil, err
	}

	return &dto.S3ApplicationMutationOutcome{App: out.App, Errors: nil}, nil
}

func (service *ApplicationService) Update(ctx context.Context, in *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
	out, err := service.bundle.appBundle.Service.Update(ctx, &dto2.ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: in.IsActive,
	})

	if nil != err {
		if err == util.ErrorUselessInput {
			if nil == in.Credentials && nil == in.Policies {
				return nil, err
			}
		}
	}

	if nil == out || nil == out.App {
		return &dto.S3ApplicationMutationOutcome{App: nil, Errors: out.Errors}, nil
	}

	err = connect.Transaction(
		ctx,
		connect.ContextToDB(ctx),
		func(tx *gorm.DB) error {
			err = service.bundle.credentialService.onAppUpdate(tx, out.App, in.Credentials)
			if nil != err {
				return err
			}

			err = service.bundle.policyService.onAppUpdate(tx, out.App, in.Policies)
			if nil != err {
				return err
			}

			return nil
		},
	)

	return &dto.S3ApplicationMutationOutcome{App: out.App, Errors: out.Errors}, err
}

func (service *ApplicationService) S3UploadToken(ctx context.Context, in dto.S3UploadTokenInput) (map[string]interface{}, error) {
	// get claims from context
	claims, ok := ctx.Value(claim.ClaimsContextKey).(*claim.Payload)
	if !ok {
		return nil, util.ErrorAuthRequired
	}

	// load application
	app, err := service.bundle.appBundle.Service.Load(ctx, in.ApplicationId)
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
