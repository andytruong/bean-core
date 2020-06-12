package model

import (
	"time"

	"bean/pkg/util/api"
)

// Any entity may have its own configuration: system, namespace, user, content, …
// Each config bucket can be configured for private or public access.
type ConfigBucket struct {
	Id          string         `json:"id"`
	Version     string         `json:"version"`
	Slug        string         `json:"slug"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	Access      api.AccessMode `json:"access"`
	HostId      string         `json:"hostId"`
	Schema      string         `json:"schema"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	IsPublished bool           `json:"isPublished"`
}