package s3

import (
	"context"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

func newResolvers(this *S3IntegrationBundle) map[string]interface{} {
	return map[string]interface{}{
		"Mutation": map[string]interface{}{
			"S3ApplicationCreate": func(ctx context.Context, input *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
				return this.AppService.Create(ctx, input)
			},
			"S3ApplicationUpdate": func(ctx context.Context, input *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
				return this.AppService.Update(ctx, input)
			},
			"S3UploadToken": func(ctx context.Context, input dto.S3UploadTokenInput) (map[string]interface{}, error) {
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
