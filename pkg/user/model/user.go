package model

import "time"

type User struct {
	ID        string    `json:"id"`
	AvatarURI *string   `json:"avatarUri"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
