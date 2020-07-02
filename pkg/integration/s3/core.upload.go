package s3

import (
	"context"
	"time"

	"github.com/minio/minio-go/v6"

	"bean/pkg/integration/s3/model"
)

type coreUpload struct {
}

func (this *coreUpload) sign(ctx context.Context, app *model.Application, bucket string, filePath string) error {
	var err error

	policy := minio.NewPostPolicy()

	if err = policy.SetBucket(bucket); nil != err {
		return err
	}

	if err = policy.SetKey(filePath); nil != err {
		return err
	}

	if err = policy.SetExpires(time.Now().UTC().Add(4 * time.Hour)); nil != err {
		return err
	}

	if err = policy.SetContentType("image/png"); nil != err {
		return err
	}

	if err = policy.SetUserMetadata("app", app.ID); nil != err {
		return err
	}

	// tag -> appId
	// tag -> sessionId -> user/device/…
	// tag -> namespaceId

	// TODO: generate per application's policy
	if err = policy.SetContentLengthRange(1, 1024*1024); nil != err {
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
