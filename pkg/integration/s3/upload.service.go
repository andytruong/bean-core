package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/app"
	"bean/pkg/integration/s3/model/dto"
)

type uploadService struct {
	bundle *Bundle
}

// TODO: move to other service
func (srv *uploadService) CreateUploadToken(ctx context.Context, in dto.UploadTokenInput) (map[string]interface{}, error) {
	// get claims from context
	claims := claim.ContextToPayload(ctx)
	if nil == claims {
		return nil, util.ErrorAuthRequired
	}

	// load application
	application, err := srv.bundle.appBundle.Service.Load(ctx, in.ApplicationId)
	if nil != err {
		return nil, err
	} else if !application.IsActive {
		return nil, app.ErrorInactiveApplication
	}

	// check upload policy
	pol, err := srv.bundle.configSrv.loadUploadPolicy(ctx, in.ApplicationId)
	if nil != err && err != gorm.ErrRecordNotFound {
		return nil, err
	} else if nil != pol {
		// TODO: pol.FileExtensions
		// TODO: file size policy

		// validate rate limit policy
		for _, limit := range pol.RateLimit {
			switch limit.Object {
			case "user":
				srv.bundle.lgr.Info("TODO.user", zap.Any("what", limit))
				fmt.Println("", limit)

			case "space":
				srv.bundle.lgr.Info("TODO.space", zap.Any("what", limit))
			}
		}
	}

	// load s3 credentials
	cre, err := srv.bundle.configSrv.loadCredentials(ctx, in.ApplicationId)
	if nil != err {
		return nil, errors.Wrap(err, "credentials not found")
	}

	client, err := srv.bundle.configSrv.client(cre)
	if nil != err {
		return nil, err
	}

	postPolicy := minio.NewPostPolicy()

	err = scalar.NoError(
		postPolicy.SetBucket(cre.Bucket),
		postPolicy.SetKey(string(in.FilePath)),
		postPolicy.SetExpires(time.Now().UTC().Add(4*time.Hour)),
		postPolicy.SetContentType(string(in.ContentType)),
		postPolicy.SetUserMetadata("app", application.ID),
		postPolicy.SetUserMetadata("sid", claims.SessionId()),
		postPolicy.SetUserMetadata("nid", claims.SpaceId()),
		postPolicy.SetContentLengthRange(1, 10*1024*1024),
	)

	if nil != err {
		return nil, err
	}

	_, formData, err := client.PresignedPostPolicy(ctx, postPolicy)
	if nil != err {
		return nil, err
	}

	response := map[string]interface{}{}
	for k, v := range formData {
		response[k] = v
	}

	return response, nil
}
