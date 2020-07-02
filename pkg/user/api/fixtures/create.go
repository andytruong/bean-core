package fixtures

import (
	"bean/components/scalar"
	"bean/pkg/user/model/dto"
	"bean/pkg/util/password"
)

func NewUserCreateInputFixture() *dto.UserCreateInput {
	passAlg, _ := password.Get("bcrypt")
	pass, _ := passAlg.Encrypt("xxxxx")

	return &dto.UserCreateInput{
		Name: &dto.UserNameInput{
			FirstName:     scalar.NilString("John"),
			LastName:      scalar.NilString("Doe"),
			PreferredName: scalar.NilString("Jon"),
		},
		Emails: &dto.UserEmailsInput{
			Primary: &dto.UserEmailInput{
				Verified: true,
				Value:    "john.doe@qa.com",
				IsActive: false,
			},
			Secondary: []*dto.UserEmailInput{
				{
					Verified: true,
					Value:    "john@doe.qa",
					IsActive: true,
				},
				{
					Verified: false,
					Value:    "john.doe@internet.qa",
					IsActive: true,
				},
			},
		},
		Password: &dto.UserPasswordInput{
			Algorithm:   passAlg.Name(),
			HashedValue: pass,
		},
		AvatarURI: scalar.NilUri("https://foo.bar"),
		IsActive:  true,
	}
}
