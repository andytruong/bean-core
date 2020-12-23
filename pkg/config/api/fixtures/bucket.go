package fixtures

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config/model/dto"
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
