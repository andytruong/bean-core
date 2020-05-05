package dto

import (
	"bean/pkg/user/model"
	"bean/pkg/util"
)

type UserCreateInput struct {
	Name      *UserNameInput   `json:"name"`
	Emails    *UserEmailsInput `json:"emails"`
	AvatarURI *string          `json:"avatarUri"`
	IsActive  bool             `json:"isActive"`
}

type UserEmailsInput struct {
	Primary   *UserEmailInput   `json:"primary"`
	Secondary []*UserEmailInput `json:"secondary"`
}

type UserEmailInput struct {
	Verified bool   `json:"verified"`
	Value    string `json:"value"`
	IsActive bool   `json:"isActive"`
}

type UserNameInput struct {
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	PreferredName *string `json:"preferredName"`
}

type UserCreateOutcome struct {
	User   *model.User   `json:"user"`
	Errors []*util.Error `json:"errors"`
}
