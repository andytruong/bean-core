package dto

import "bean/pkg/util/api"

type BucketUpdateInput struct {
	Id          string
	Version     string
	Title       *string         `json:"title"`
	Description *string         `json:"description"`
	Access      *api.AccessMode `json:"access"`
	Schema      *string         `json:"schema"`
	IsPublished bool            `json:"isPublished"`
}
