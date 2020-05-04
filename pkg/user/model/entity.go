package model

import "time"

type User struct {
	ID        string      `json:"id"`
	Name      *UserName   `json:"name"`
	Emails    *UserEmails `json:"emails"`
	AvatarURI *string     `json:"avatarUri"`
	IsActive  bool        `json:"isActive"`
}

type UserEmail struct {
	ID        string    `json:"id"`
	Verified  bool      `json:"verified"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsActive  bool      `json:"isActive"`
}

type UserEmailInput struct {
	Verified bool   `json:"verified"`
	Value    string `json:"value"`
	IsActive bool   `json:"isActive"`
}

type UserEmails struct {
	Primary   *UserEmail   `json:"primary"`
	Secondary []*UserEmail `json:"secondary"`
}

type UserName struct {
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	PrefferedName *string `json:"prefferedName"`
}
