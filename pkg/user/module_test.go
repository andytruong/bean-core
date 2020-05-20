package user

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"bean/pkg/user/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/password"
)

func TestUserMutationResolver_UserCreate(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	mod := NewUserModule(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(mod, db)

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
		out, err := mod.UserCreate(context.Background(), input)
		ass.NoError(err)
		ass.Empty(out.Errors)
		ass.Equal("https://foo.bar", string(*out.User.AvatarURI))

		theUser, err := mod.User(context.Background(), out.User.ID)
		ass.NoError(err)
		ass.True(theUser.CreatedAt.UnixNano() >= now.UnixNano())
		ass.Equal("https://foo.bar", string(*theUser.AvatarURI))
	})

	t.Run("error by email duplication", func(t *testing.T) {
		ass.True(true, "TODO")
	})

	t.Run("email is casted to lower-case", func(t *testing.T) {
		ass.True(true, "TODO")
	})
}
