package model

import (
	"time"
)

type Policy struct {
	ID            string     `json:"id"`
	ApplicationId string     `json:"applicationId"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Kind          PolicyKind `json:"kind"`
	Value         string     `json:"value"`
}

func (Policy) TableName() string {
	return "s3_application_policy"
}

type PolicyKind string

const (
	// Example: "pdf txt zip gz"
	PolicyKindFileExtensions PolicyKind = "file_extensions"

	// Example: "1MB/user/hour", "1GB/space/hour"
	PolicyKindRateLimit PolicyKind = "rate_limit"
)
