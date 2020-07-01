package dto

import (
	"bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/api/scalar"
)

type UserCreateInput struct {
	Name      *UserNameInput     `json:"name"`
	Emails    *UserEmailsInput   `json:"emails"`
	Password  *UserPasswordInput `json:"password"`
	AvatarURI *scalar.Uri        `json:"avatarUri"`
	IsActive  bool               `json:"isActive"`
}

type UserEmailsInput struct {
	Primary   *UserEmailInput   `json:"primary"`
	Secondary []*UserEmailInput `json:"secondary"`
}

type UserEmailInput struct {
	Verified bool                `json:"verified"`
	Value    scalar.EmailAddress `json:"value"`
	IsActive bool                `json:"isActive"`
}

type UserNameInput struct {
	FirstName     *string `json:"firstName"`
	LastName      *string `json:"lastName"`
	PreferredName *string `json:"preferredName"`
}

type UserPasswordInput struct {
	Algorithm   string `json:"algorithm"`
	HashedValue string `json:"hashedValue"`
}

type UserMutationOutcome struct {
	Errors []*util.Error `json:"errors"`
	User   *model.User   `json:"user"`
}
