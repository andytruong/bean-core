package model

import (
	"bean/components/scalar"
)

type Credentials struct {
	ID            string     `json:"id"`
	ApplicationId string     `json:"applicationId"`
	Endpoint      scalar.Uri `json:"endpoint"`
	Bucket        string     `json:"bucket"`
	IsSecure      bool       `json:"isSecure"`
	AccessKey     string     `json:"accessKey"`
	SecretKey     string     `json:"secretKey"`
}

func (Credentials) TableName() string {
	return "s3_credentials"
}
