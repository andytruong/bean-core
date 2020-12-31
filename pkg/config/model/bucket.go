package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/qri-io/jsonschema"

	"bean/components/scalar"
)

// Any entity may have its own configuration: system, space, user, content, …
// Each config bucket can be configured for private or public access.
type ConfigBucket struct {
	Id          string            `json:"id"`
	Version     string            `json:"version"`
	Slug        string            `json:"slug"`
	Title       string            `json:"title"`
	Description *string           `json:"description"`
	Access      scalar.AccessMode `json:"access"`
	HostId      string            `json:"hostId"`
	Schema      string            `json:"schema"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	IsPublished bool              `json:"isPublished"`
}

func (ConfigBucket) TableName() string {
	return "config_buckets"
}

func (bucket ConfigBucket) Validate(ctx context.Context, value string) ([]string, error) {
	var (
		err     error
		reasons []string
	)

	rs := &jsonschema.Schema{}
	err = json.Unmarshal([]byte(bucket.Schema), rs)
	if nil != err {
		return nil, err
	}

	explanations, err := rs.ValidateBytes(ctx, []byte(value))
	if nil != err {
		return nil, err
	}

	if len(explanations) > 0 {
		reasons = []string{}
		for _, reason := range explanations {
			reasons = append(reasons, reason.Error())
		}
	}

	return reasons, nil
}
