package model

import (
	"time"

	"bean/pkg/util/api"
)

type UserEmails struct {
	Primary   *UserEmail   `json:"primary"`
	Secondary []*UserEmail `json:"secondary"`
}

type UserEmail struct {
	ID         string           `json:"id"`
	UserId     string           `json:"userId"`
	Value      api.EmailAddress `json:"value"`
	IsActive   bool             `json:"isActive"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
	IsPrimary  bool             `json:"isPrimary"`
	IsVerified bool             `gorm:"-"`
}

type UserUnverifiedEmail struct {
	ID        string           `json:"id"        gorm:"primary_key"`
	Value     api.EmailAddress `json:"value"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}
