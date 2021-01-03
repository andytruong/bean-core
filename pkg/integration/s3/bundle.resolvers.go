package s3

import (
	"context"

	appModel "bean/pkg/app/model"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

func newResolvers(bundle *Bundle) map[string]interface{} {
	return map[string]interface{}{
		"Application": map[string]interface{}{
			"S3Credentials": func(ctx context.Context, app *appModel.Application) (*model.S3Credentials, error) {
				return bundle.configSrv.loadCredentials(ctx, app.ID)
			},
			"S3UploadPolices": func(ctx context.Context, app *appModel.Application) (*model.S3UploadPolicy, error) {
				return bundle.configSrv.loadUploadPolicy(ctx, app.ID)
			},
		},
		"Mutation": map[string]interface{}{
			"S3Mutation": func(ctx context.Context) (*dto.S3Mutation, error) {
				return &dto.S3Mutation{}, nil
			},
		},
		"S3Mutation": map[string]interface{}{
			"SaveCredentials": func(ctx context.Context, _ *dto.S3Mutation, in dto.S3CredentialsInput) (*dto.S3CredentialsOutcome, error) {
				cre, err := bundle.configSrv.saveCredentials(ctx, in)
				if nil != err {
					return nil, err
				}

				return &dto.S3CredentialsOutcome{Errors: nil, Credentials: cre}, err
			},
			"SaveUploadPolicies": func(ctx context.Context, _ *dto.S3Mutation, in dto.UploadPolicyInput) (*dto.S3UploadPolicyOutcome, error) {
				policy, err := bundle.configSrv.saveUploadPolicy(ctx, in)
				if nil != err {
					return nil, err
				}

				return &dto.S3UploadPolicyOutcome{Errors: nil, Policy: policy}, nil
			},
			"UploadToken": func(ctx context.Context, _ *dto.S3Mutation, in dto.UploadTokenInput) (map[string]interface{}, error) {
				return bundle.uploadSrv.CreateUploadToken(ctx, in)
			},
		},
	}
}
