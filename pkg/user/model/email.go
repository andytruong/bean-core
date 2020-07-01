package model

import (
	"time"

	"bean/pkg/util/api/scalar"
)

type UserEmails struct {
	Primary   *UserEmail   `json:"primary"`
	Secondary []*UserEmail `json:"secondary"`
}

type UserEmail struct {
	ID         string              `json:"id"`
	UserId     string              `json:"userId"`
	Value      scalar.EmailAddress `json:"value"`
	IsActive   bool                `json:"isActive"`
	CreatedAt  time.Time           `json:"createdAt"`
	UpdatedAt  time.Time           `json:"updatedAt"`
	IsPrimary  bool                `json:"isPrimary"`
	IsVerified bool                `gorm:"-"`
}

type UserUnverifiedEmail struct {
	ID        string              `json:"id"        gorm:"primary_key"`
	Value     scalar.EmailAddress `json:"value"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}
