package dto

import "bean/pkg/util"

type S3ApplicationUpdateInput struct {
	Id          string                               `json:"id"`
	Version     string                               `json:"version"`
	IsActive    *bool                                `json:"isActive"`
	Slug        *string                              `json:"slug"`
	Credentials *S3ApplicationCredentialsUpdateInput `json:"credentials"`
}

type S3ApplicationCredentialsUpdateInput struct {
	Endpoint  *util.Uri `json:"endpoint"`
	IsSecure  *bool     `json:"isSecure"`
	AccessKey *string   `json:"accessKey"`
	SecretKey *string   `json:"secretKey"`
}
