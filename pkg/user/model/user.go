package model

import (
	"time"

	"bean/components/scalar"
	"bean/pkg/util/api"
)

type User struct {
	ID        string       `json:"id"`
	Version   string       `json:"version"`
	IsActive  bool         `json:"isActive"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	AvatarURI *scalar.Uri  `json:"avatarUri"`
	Language  api.Language `json:"language"`
}
