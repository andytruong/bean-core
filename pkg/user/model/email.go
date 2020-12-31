package model

import (
	"time"

	"bean/components/connect"
	"bean/components/scalar"
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

func (UserEmail) Name() string {
	return connect.TableUserEmail
}

type UserUnverifiedEmail struct {
	ID        string              `json:"id"        gorm:"primary_key"`
	Value     scalar.EmailAddress `json:"value"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

func (UserUnverifiedEmail) Name() string {
	return connect.TableUserEmailUnverified
}
