package dto

import (
	"bean/components/scalar"
	util2 "bean/components/util"
	"bean/pkg/integration/s3/model"
)

type S3ApplicationCreateInput struct {
	IsActive    bool                                `json:"isActive"`
	Credentials S3ApplicationCredentialsCreateInput `json:"credentials"`
	Policies    []S3ApplicationPolicyCreateInput    `json:"policies"`
}

type S3ApplicationCredentialsCreateInput struct {
	Endpoint  scalar.Uri `json:"endpoint"`
	Bucket    string     `json:"bucket"`
	IsSecure  bool       `json:"isSecure"`
	AccessKey string     `json:"accessKey"`
	SecretKey string     `json:"secretKey"`
}

type S3ApplicationMutationOutcome struct {
	App    *model.Application `json:"application"`
	Errors []*util2.Error     `json:"errors"`
}
