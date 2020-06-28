package model

import (
	"time"

	"bean/pkg/util"
)

type Application struct {
	Slug      string     `json:"name"`
	ID        string     `json:"id"`
	Version   string     `json:"version"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type Credentials struct {
	ID               string   `json:"id"`
	Version          string   `json:"version"`
	ApplicationId    string   `json:"applicationId"`
	Endpoint         util.Uri `json:"endpoint"`
	EncryptedKeyPair string   `json:"encryptedKeyPair"`
	IsSecure         bool     `json:"isSecure"`
}

type Policy struct {
	ID            string     `json:"id"`
	Version       string     `json:"version"`
	ApplicationId string     `json:"applicationId"`
	IsActive      bool       `json:"isActive"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	Kind          PolicyKind `json:"kind"`
	Value         string     `json:"value"`
}

type PolicyKind string

const (
	PolicyKindFileExtensions PolicyKind = "file_extensions"
	PolicyKindRateLimit      PolicyKind = "rate_limit"
)
