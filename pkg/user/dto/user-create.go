package dto

import (
	"bean/pkg/user/model"
	"bean/pkg/util"
)

type UserCreateInput struct {
	ID        *string          `json:"id"`
	Name      *UserNameInput   `json:"name"`
	Emails    *UserEmailsInput `json:"emails"`
	AvatarURI *string          `json:"avatarUri"`
	IsActive  bool             `json:"isActive"`
}

type UserEmailsInput struct {
	Primary   *model.UserEmailInput   `json:"primary"`
	Secondary []*model.UserEmailInput `json:"secondary"`
}

type UserNameInput struct {
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	PrefferedName *string `json:"prefferedName"`
}

type UserCreateOutcome struct {
	User   *model.User   `json:"user"`
	Errors []*util.Error `json:"errors"`
}
