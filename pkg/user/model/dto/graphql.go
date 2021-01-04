package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/user/model"
)

type UserQuery struct{}

// Mutation
type (
	UserMutation struct{}

	UserMutationOutcome struct {
		Errors []*util.Error `json:"errors"`
		User   *model.User   `json:"user"`
	}

	// create
	UserCreateInput struct {
		Name      *UserNameInput     `json:"name"`
		Emails    *UserEmailsInput   `json:"emails"`
		Password  *UserPasswordInput `json:"password"`
		AvatarURI *scalar.Uri        `json:"avatarUri"`
		IsActive  bool               `json:"isActive"`
	}

	UserEmailsInput struct {
		Primary   *UserEmailInput   `json:"primary"`
		Secondary []*UserEmailInput `json:"secondary"`
	}

	UserEmailInput struct {
		Verified bool                `json:"verified"`
		Value    scalar.EmailAddress `json:"value"`
		IsActive bool                `json:"isActive"`
	}

	UserNameInput struct {
		FirstName     *string `json:"firstName"`
		LastName      *string `json:"lastName"`
		PreferredName *string `json:"preferredName"`
	}

	UserPasswordInput struct {
		HashedValue string `json:"hashedValue"`
	}

	// Update
	UserUpdateInput struct {
		ID      string                 `json:"id"`
		Version string                 `json:"version"`
		Values  *UserUpdateValuesInput `json:"values"`
	}

	UserUpdateValuesInput struct {
		Password *UserPasswordInput `json:"password"`
	}
)
