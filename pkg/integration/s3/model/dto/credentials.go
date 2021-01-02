package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	appModel "bean/pkg/app/model"
)

type S3ApplicationCredentialsInput struct {
	Version       string     `json:"version"`
	ApplicationId string     `json:"applicationId"`
	Endpoint      scalar.Uri `json:"endpoint"`
	Bucket        string     `json:"bucket"`
	IsSecure      bool       `json:"isSecure"`
	AccessKey     string     `json:"accessKey"`
	SecretKey     string     `json:"secretKey"`
}

// TODO: Remove
type S3ApplicationCredentialsUpdateInput struct {
	Endpoint  *scalar.Uri `json:"endpoint"`
	Bucket    *string     `json:"bucket"`
	IsSecure  *bool       `json:"isSecure"`
	AccessKey *string     `json:"accessKey"`
	SecretKey *string     `json:"secretKey"`
}

type S3ApplicationMutationOutcome struct {
	App    *appModel.Application `json:"application"`
	Errors []*util.Error         `json:"errors"`
}
