package model

import (
	"time"
	
	"bean/components/scalar"
)

type Application struct {
	ID        string     `json:"id"`
	Version   string     `json:"version"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (this Application) TableName() string {
	return "s3_application"
}

type Credentials struct {
	ID            string     `json:"id"`
	ApplicationId string     `json:"applicationId"`
	Endpoint      scalar.Uri `json:"endpoint"`
	Bucket        string     `json:"bucket"`
	IsSecure      bool       `json:"isSecure"`
	AccessKey     string     `json:"accessKey"`
	SecretKey     string     `json:"secretKey"`
}

func (this Credentials) TableName() string {
	return "s3_credentials"
}

type Policy struct {
	ID            string     `json:"id"`
	ApplicationId string     `json:"applicationId"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Kind          PolicyKind `json:"kind"`
	Value         string     `json:"value"`
}

func (this Policy) TableName() string {
	return "s3_application_policy"
}

type PolicyKind string

const (
	// Example: "pdf txt zip gz"
	PolicyKindFileExtensions PolicyKind = "file_extensions"
	
	// Example: "1MB/user/hour", "1GB/space/hour"
	PolicyKindRateLimit PolicyKind = "rate_limit"
)
