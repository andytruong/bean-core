package s3

import (
	"context"
	
	appModel "bean/pkg/app/model"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

func newResolvers(this *S3Bundle) map[string]interface{} {
	return map[string]interface{}{
		"Application": map[string]interface{}{
			"Credentials": func(ctx context.Context, obj *appModel.Application) (*model.Credentials, error) {
				return this.credentialService.load(ctx, obj.ID)
			},
			"Polices": func(ctx context.Context, obj *appModel.Application) ([]*model.Policy, error) {
				return this.policyService.loadByApplicationId(ctx, obj.ID)
			},
		},
		"Mutation": map[string]interface{}{
			"S3Mutation": func(ctx context.Context) (*dto.S3Mutation, error) {
				return &dto.S3Mutation{}, nil
			},
		},
		"S3Mutation": map[string]interface{}{
			"Application": func(ctx context.Context) (*dto.S3ApplicationMutation, error) {
				return &dto.S3ApplicationMutation{}, nil
			},
			"Upload": func(ctx context.Context, _ *dto.S3Mutation) (*dto.S3UploadMutation, error) {
				return &dto.S3UploadMutation{}, nil
			},
		},
		"S3UploadMutation": map[string]interface{}{
			"Token": func(ctx context.Context, _ *dto.S3UploadMutation, input dto.S3UploadTokenInput) (map[string]interface{}, error) {
				return this.AppService.S3UploadToken(ctx, input)
			},
		},
	}
}
