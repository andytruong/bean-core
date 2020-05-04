package model

import "time"

type UserEmails struct {
	Primary   *UserEmail   `json:"primary"`
	Secondary []*UserEmail `json:"secondary"`
}

type UserEmail struct {
	ID        string    `json:"id"`
	UserId    string    `json:"userId"`
	Verified  bool      `json:"verified"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
}
