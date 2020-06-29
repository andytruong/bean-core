package dto

import (
	"bean/pkg/integration/s3/model"
	"bean/pkg/util"
)

type S3ApplicationCreateInput struct {
	Slug        string                              `json:"slug"`
	IsActive    bool                                `json:"isActive"`
	Credentials S3ApplicationCredentialsCreateInput `json:"credentials"`
	Polices     []S3ApplicationPolicyCreateInput    `json:"policies"`
}

type S3ApplicationCredentialsCreateInput struct {
	Endpoint  util.Uri `json:"endpoint"`
	IsSecure  bool     `json:"isSecure"`
	AccessKey string   `json:"accessKey"`
	SecretKey string   `json:"secretKey"`
}

type S3ApplicationMutationOutcome struct {
	App    *model.Application `json:"application"`
	Errors []*util.Error      `json:"errors"`
}
