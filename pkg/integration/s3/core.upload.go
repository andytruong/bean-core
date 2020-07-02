package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v6"

	"bean/components/scalar"
	"bean/pkg/integration/s3/model"
	"bean/pkg/util"
)

type coreUpload struct {
}

func (this *coreUpload) sign(
	ctx context.Context,
	claims util.Claims,
	app *model.Application,
	bucket string,
	filePath string,
) error {
	policy := minio.NewPostPolicy()
	
	err := scalar.NoError(
		policy.SetBucket(bucket),
		policy.SetKey(filePath),
		policy.SetExpires(time.Now().UTC().Add(4*time.Hour)),
		policy.SetContentType("image/png"),
		policy.SetUserMetadata("app", app.ID),
		policy.SetUserMetadata("sid", claims.SessionId()),
		policy.SetUserMetadata("nid", claims.NamespaceId()),
		policy.SetContentLengthRange(1, 18*1024*1024), // TODO: generate per application's policy
	)

	if nil != err {
		return err
	}

	// bucket?
	// actor?
	// namespace?
	// file: path, size, tags?
	// flood detect
	// error code

	panic("wip")
}
