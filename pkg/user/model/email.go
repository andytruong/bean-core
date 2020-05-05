package model

import (
	"time"

	"bean/pkg/util"
)

type UserEmails struct {
	Primary   *UserEmail   `json:"primary"`
	Secondary []*UserEmail `json:"secondary"`
}

type UserEmail struct {
	ID        string            `json:"id"`
	UserId    string            `json:"userId"`
	Verified  bool              `json:"verified"`
	Value     util.EmailAddress `json:"value"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	IsActive  bool              `json:"isActive"`
}
