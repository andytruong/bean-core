package dto

import (
	"bean/components/scalar"
)

type S3ApplicationUpdateInput struct {
	Id          string                               `json:"id"`
	Version     string                               `json:"version"`
	IsActive    *bool                                `json:"isActive"`
	Credentials *S3ApplicationCredentialsUpdateInput `json:"credentials"`
	Policies    *S3ApplicationPolicyMutationInput    `json:"policies"`
}

type S3ApplicationCredentialsUpdateInput struct {
	Endpoint  *scalar.Uri `json:"endpoint"`
	Bucket    *string     `json:"bucket"`
	IsSecure  *bool       `json:"isSecure"`
	AccessKey *string     `json:"accessKey"`
	SecretKey *string     `json:"secretKey"`
}