package s3

import (
	"context"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type ApplicationResolver struct {
	bean *S3IntegrationBean
}

func (this *ApplicationResolver) S3ApplicationCreate(ctx context.Context, input *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	return this.bean.CoreApp.Create(ctx, input)
}

func (this *ApplicationResolver) Polices(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
	return this.bean.corePolicy.loadByApplicationId(ctx, obj.ID)
}

func (this *ApplicationResolver) Credentials(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
	return this.bean.coreCredentials.loadByApplicationId(ctx, obj.ID)
}
