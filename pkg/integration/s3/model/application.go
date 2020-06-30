package model

import (
	"time"

	"bean/pkg/util"
)

type Application struct {
	ID        string     `json:"id"`
	Version   string     `json:"version"`
	Slug      string     `json:"name"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type Credentials struct {
	ID            string   `json:"id"`
	ApplicationId string   `json:"applicationId"`
	Endpoint      util.Uri `json:"endpoint"`
	IsSecure      bool     `json:"isSecure"`
	AccessKey     string   `json:"accessKey"`
	SecretKey     string   `json:"secretKey"`
}

type Policy struct {
	ID            string     `json:"id"`
	ApplicationId string     `json:"applicationId"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Kind          PolicyKind `json:"kind"`
	Value         string     `json:"value"`
}

type PolicyKind string

const (
	// Example: "pdf txt zip gz"
	PolicyKindFileExtensions PolicyKind = "file_extensions"

	// Example: "1MB/user/hour", "1GB/namespace/hour"
	PolicyKindRateLimit PolicyKind = "rate_limit"
)
