package model

import (
	"time"

	"bean/pkg/util"
)

type User struct {
	ID        string    `json:"id"`
	Version   string    `json:"version"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	AvatarURI *util.Uri `json:"avatarUri"`
}
