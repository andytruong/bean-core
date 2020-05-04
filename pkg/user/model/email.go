package model

import "time"

type UserEmail struct {
	ID        string    `json:"id"`
	Verified  bool      `json:"verified"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
}

type UserEmails struct {
	Primary   *UserEmail   `json:"primary"`
	Secondary []*UserEmail `json:"secondary"`
}
