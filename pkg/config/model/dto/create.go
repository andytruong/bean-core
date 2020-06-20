package dto

import (
	"bean/pkg/config/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type BucketCreateInput struct {
	// ID of host entity.
	//  If bucket is created for namespace -> bucket.HostId -> namespace.ID
	//  If bucket is created for user      -> bucket.HostId -> user.ID
	HostId      string          `json:"hostId"`
	Slug        *string         `json:"slug"`
	Title       *string         `json:"title"`
	Description *string         `json:"description"`
	Access      *api.AccessMode `json:"access"`
	Schema      string          `json:"schema"`
	IsPublished bool            `json:"isPublished"`
}

type BucketMutationOutcome struct {
	Errors []util.Error        `json:"errors"`
	Bucket *model.ConfigBucket `json:"bucket"`
}
