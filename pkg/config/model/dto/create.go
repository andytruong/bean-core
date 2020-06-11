package dto

import (
	"bean/pkg/config/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type BucketCreateInput struct {
	HostId      string          `json:"hostId"`
	Slug        *string         `json:"slug"`
	Title       *string         `json:"title"`
	Description *string         `json:"description"`
	Access      *api.AccessMode `json:"access"`
	Schema      string          `json:"schema"`
}

type BucketMutationOutcome struct {
	Errors []util.Error        `json:"errors"`
	Bucket *model.ConfigBucket `json:"bucket"`
}
