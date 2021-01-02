package model

import (
	"bean/components/scalar"
)

type Credentials struct {
	Id        string     `json:"id,omitempty"`
	Version   string     `json:"version,omitempty"`
	Endpoint  scalar.Uri `json:"endpoint"`
	Bucket    string     `json:"bucket"`
	IsSecure  bool       `json:"isSecure"`
	AccessKey string     `json:"accessKey"`
	SecretKey string     `json:"secretKey"`
}
