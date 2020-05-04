package model

type User struct {
	ID        string      `json:"id"`
	Name      *UserName   `json:"name"`
	Emails    *UserEmails `json:"emails"`
	AvatarURI *string     `json:"avatarUri"`
	IsActive  bool        `json:"isActive"`
}
