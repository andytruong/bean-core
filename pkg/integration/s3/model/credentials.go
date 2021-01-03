package model

import (
	"bytes"
	"encoding/json"

	"bean/components/scalar"
	"bean/pkg/infra/api"
)

type S3Credentials struct {
	Id        string     `json:"id,omitempty"`
	Version   string     `json:"version,omitempty"`
	Endpoint  scalar.Uri `json:"endpoint"`
	Bucket    string     `json:"bucket"`
	IsSecure  bool       `json:"isSecure"`
	AccessKey string     `json:"accessKey"`
	SecretKey string     `json:"secretKey"`
}

type S3UploadPolicy struct {
	Id             string                  `json:"id,omitempty"`
	Version        string                  `json:"version,omitempty"`
	FileExtensions []api.FileType          `json:"fileExtensions"`
	RateLimit      []UploadRateLimitPolicy `json:"rateLimit"`
}

func (pl S3UploadPolicy) EqualTo(newVersion S3UploadPolicy) (bool, error) {
	pl.Id = ""
	pl.Version = ""
	newVersion.Id = ""
	newVersion.Version = ""

	oldBytes, _ := json.Marshal(pl)
	newBytes, _ := json.Marshal(newVersion)

	return bytes.Equal(oldBytes, newBytes), nil
}

type UploadRateLimitPolicy struct {
	Value    string `json:"value"`    // example: 1MB, 1GB
	Object   string `json:"object"`   // oneOf: user, space
	Interval string `json:"interval"` // oneOf: minute/hour/day
}
