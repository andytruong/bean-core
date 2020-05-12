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
	Value     util.EmailAddress `json:"value"`
	IsActive  bool              `json:"isActive"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Verified  bool              `json:"isVerified"`
}

type UserUnverifiedEmail struct {
	ID        string            `json:"id"        gorm:"primary_key"`
	Value     util.EmailAddress `json:"value"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}
