package s3

import (
	"context"
	
	"bean/pkg/access/model/dto"
	"bean/pkg/integration/s3/model"
	dto3 "bean/pkg/integration/s3/model/dto"
)

func newGraphqlResolver(bundle *S3IntegrationBundle) map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{},
		"Mutation": map[string]interface{}{
			"S3ApplicationCreate": func(ctx context.Context, input *dto3.S3ApplicationCreateInput) (*dto3.S3ApplicationMutationOutcome, error) {
				panic("wip")
			},
			"SessionArchive": func(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
				panic("wip")
			},
		},
		"Application": map[string]interface{}{
			"Polices": func(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
				return bundle.policyService.loadByApplicationId(ctx, obj.ID)
			},
			"Credentials": func(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
				return bundle.credentialService.loadByApplicationId(ctx, obj.ID)
			},
		},
	}
}
