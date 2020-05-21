package fixtures

import (
	"bean/pkg/user/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/password"
)

func NewUserCreateInputFixture() *dto.UserCreateInput {
	passAlg, _ := password.Get("bcrypt")
	pass, _ := passAlg.Encrypt("xxxxx")

	return &dto.UserCreateInput{
		Name: &dto.UserNameInput{
			FirstName:     util.NilString("John"),
			LastName:      util.NilString("Doe"),
			PreferredName: util.NilString("Jon"),
		},
		Emails: &dto.UserEmailsInput{
			Primary: &dto.UserEmailInput{
				Verified: true,
				Value:    "john.doe@qa.com",
				IsActive: false,
			},
			Secondary: []*dto.UserEmailInput{
				{
					Verified: false,
					Value:    "john@doe.qa",
					IsActive: true,
				},
				{
					Verified: true,
					Value:    "john.doe@internet.qa",
					IsActive: true,
				},
			},
		},
		Password: &dto.UserPasswordInput{
			Algorithm:   passAlg.Name(),
			HashedValue: pass,
		},
		AvatarURI: util.NilUri("https://foo.bar"),
		IsActive:  true,
	}
}
