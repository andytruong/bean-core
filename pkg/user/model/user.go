package model

import (
	"time"

	"bean/pkg/util"
)

type User struct {
	ID        string    `json:"id"`
	AvatarURI *util.Uri `json:"avatarUri"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
