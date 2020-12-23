package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	model2 "bean/pkg/app/model"
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
	App    *model2.Application `json:"application"`
	Errors []*util.Error       `json:"errors"`
}
