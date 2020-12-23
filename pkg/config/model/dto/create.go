package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config/model"
)

type BucketCreateInput struct {
	// ID of host entity.
	//  If bucket is created for space -> bucket.HostId -> space.ID
	//  If bucket is created for user      -> bucket.HostId -> user.ID
	HostId      string             `json:"hostId"`
	Slug        *string            `json:"slug"`
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Access      *scalar.AccessMode `json:"access"`
	Schema      string             `json:"schema"`
	IsPublished bool               `json:"isPublished"`
}

type BucketMutationOutcome struct {
	Errors []util.Error        `json:"errors"`
	Bucket *model.ConfigBucket `json:"bucket"`
}
