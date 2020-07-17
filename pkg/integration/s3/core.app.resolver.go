package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
)

type ApplicationResolver struct {
	bean *S3IntegrationBean
}

func (this *ApplicationResolver) S3ApplicationCreate(ctx context.Context, input *dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	return this.bean.CoreApp.Create(ctx, input)
}

func (this *ApplicationResolver) S3ApplicationUpdate(ctx context.Context, in *dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
	return this.bean.CoreApp.Update(ctx, in)
}

func (this *ApplicationResolver) Polices(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
	return this.bean.corePolicy.loadByApplicationId(ctx, obj.ID)
}

func (this *ApplicationResolver) Credentials(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
	return this.bean.coreCredentials.loadByApplicationId(ctx, obj.ID)
}

func (this *ApplicationResolver) S3UploadToken(ctx context.Context, in dto.S3UploadTokenInput) (map[string]interface{}, error) {
	// get claims from context
	claims, ok := ctx.Value(claim.ContextKey).(*claim.Payload)
	if !ok {
		return nil, util.ErrorAuthRequired
	}

	// load application
	app, err := this.bean.CoreApp.Load(ctx, in.ApplicationId)
	if nil != err {
		return nil, err
	} else {
		cred, err := this.bean.coreCredentials.loadByApplicationId(ctx, in.ApplicationId)
		if nil != err {
			return nil, err
		} else if client, err := this.bean.coreCredentials.client(cred); nil != err {
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
				policy.SetUserMetadata("nid", claims.NamespaceId()),
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
