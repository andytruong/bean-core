package s3

import (
	"context"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

func newResolvers(this *S3Bundle) map[string]interface{} {
	return map[string]interface{}{
		"Mutation": map[string]interface{}{
			"S3Mutation": func(ctx context.Context) (*dto.S3Mutation, error) {
				return &dto.S3Mutation{}, nil
			},
		},
		"S3Mutation": map[string]interface{}{
			"Application": func(ctx context.Context) (*dto.S3ApplicationMutation, error) {
				return &dto.S3ApplicationMutation{}, nil
			},
			"Upload": func(ctx context.Context) (*dto.S3UploadMutation, error) {
				return &dto.S3UploadMutation{}, nil
			},
		},
		"S3ApplicationMutation": map[string]interface{}{
			"Create": func(ctx context.Context, input *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
				return this.AppService.Create(ctx, input)
			},
			"Update": func(ctx context.Context, input *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
				return this.AppService.Update(ctx, input)
			},
		},
		"S3UploadMutation": map[string]interface{}{
			"Token": func(ctx context.Context, input dto.S3UploadTokenInput) (map[string]interface{}, error) {
				return this.AppService.S3UploadToken(ctx, input)
			},
		},
		"Application": map[string]interface{}{
			"Polices": func(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
				return this.policyService.loadByApplicationId(ctx, obj.ID)
			},
			"Credentials": func(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
				return this.credentialService.loadByApplicationId(ctx, obj.ID)
			},
		},
	}
}
