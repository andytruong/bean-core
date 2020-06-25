package dto

import "bean/pkg/util"

type S3ApplicationCreateInput struct {
	Slug        string `json:"slug"`
	IsActive    bool   `json:"isActive"`
	Credentials S3ApplicationCredentialsInput
}

type S3ApplicationCredentialsInput struct {
	Endpoint  util.Uri `json:"endpoint"`
	IsSecure  bool     `json:"isSecure"`
	AccessKey string   `json:"accessKey"`
	SecretKey string   `json:"secretKey"`
}
