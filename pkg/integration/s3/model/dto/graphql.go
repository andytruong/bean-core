package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/infra/api"
	"bean/pkg/integration/s3/model"
)

type S3Mutation struct{}

type S3CredentialsInput struct {
	Version       string     `json:"version"`
	ApplicationId string     `json:"applicationId"`
	Endpoint      scalar.Uri `json:"endpoint"`
	Bucket        string     `json:"bucket"`
	IsSecure      bool       `json:"isSecure"`
	AccessKey     string     `json:"accessKey"`
	SecretKey     string     `json:"secretKey"`
}

type S3CredentialsOutcome struct {
	Errors      []*util.Error        `json:"errors"`
	Credentials *model.S3Credentials `json:"credentials"`
}

type UploadTokenInput struct {
	ApplicationId string             `json:"applicationId"`
	FilePath      scalar.Uri         `json:"filePath"`
	ContentType   scalar.ContentType `json:"contentType"`
}

type UploadPolicyInput struct {
	Version        string                 `json:"version"`
	ApplicationId  string                 `json:"applicationId"`
	FileExtensions []api.FileType         `json:"fileExtensions"`
	RateLimit      []UploadRateLimitInput `json:"rateLimit"`
}

type UploadRateLimitInput struct {
	Value    string `json:"value"`
	Object   string `json:"object"`
	Interval string `json:"interval"`
}

type S3UploadPolicyOutcome struct {
	Errors []*util.Error         `json:"errors"`
	Policy *model.S3UploadPolicy `json:"policy"`
}
