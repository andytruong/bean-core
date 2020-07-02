package fixtures

import (
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api/scalar"
)

func NewConfigBucketCreate(access scalar.AccessMode) *dto.BucketCreateInput {
	id := util.MockIdentifier()

	return &dto.BucketCreateInput{
		HostId:      id.MustULID(),
		Slug:        scalar.NilString("doe"),
		Title:       scalar.NilString("Doe"),
		Description: scalar.NilString("Just for John Doe"),
		Access:      &access,
		Schema:      `{"type:"number"}`,
		IsPublished: false,
	}
}
