package user

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"bean/pkg/user/dto"
	"bean/pkg/util"
	"bean/pkg/util/password"
)

func TestNewUserModule(t *testing.T) {
	ass := assert.New(t)
	_, err := NewUserModule(util.MockDatabase(), util.MockLogger(), util.MockIdentifier())
	ass.NoError(err)
}

func TestUserMutationResolver_UserCreate(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	module, err := NewUserModule(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(module, db)
	ass.NoError(err)

	passAlg, _ := password.Get("bcrypt")
	pass, _ := passAlg.Encrypt("xxxxx")
	input := &dto.UserCreateInput{
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

	t.Run("test happy case, no error", func(t *testing.T) {
		now := time.Now()
		outcome, err := module.Mutation.UserCreate(context.Background(), input)
		ass.NoError(err)
		ass.Empty(outcome.Errors)
		ass.Equal("https://foo.bar", string(*outcome.User.AvatarURI))

		theUser, err := module.Query.User(context.Background(), outcome.User.ID)
		ass.NoError(err)
		ass.True(theUser.CreatedAt.UnixNano() >= now.UnixNano())
		ass.Equal("https://foo.bar", string(*theUser.AvatarURI))
	})

	t.Run("error by email duplication", func(t *testing.T) {
		ass.True(true, "TODO")
	})
}
