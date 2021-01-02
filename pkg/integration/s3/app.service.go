package s3

import (
	"context"
	"time"
	
	"github.com/minio/minio-go/v7"
	
	"bean/components/claim"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/integration/s3/model/dto"
)

type AppService struct {
	bundle *S3Bundle
}

// TODO: move to other service
func (srv *AppService) S3UploadToken(ctx context.Context, in dto.S3UploadTokenInput) (map[string]interface{}, error) {
	// get claims from context
	claims := claim.ContextToPayload(ctx)
	if nil == claims {
		return nil, util.ErrorAuthRequired
	}
	
	// load application
	app, err := srv.bundle.appBundle.Service.Load(ctx, in.ApplicationId)
	if nil != err {
		return nil, err
	} else {
		cred, err := srv.bundle.credentialService.load(ctx, in.ApplicationId)
		if nil != err {
			return nil, err
		} else if client, err := srv.bundle.credentialService.client(cred); nil != err {
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
