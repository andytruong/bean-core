package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"bean/pkg/user/dto"
	"bean/pkg/util"
	"bean/pkg/util/password"
)

func TestNewUserModule(t *testing.T) {
	// container := infra.NewMockContainer
}

func TestUserMutationResolver_UserCreate(t *testing.T) {
	ass := assert.New(t)
	db, err := gorm.Open("sqlite3", ":memory:")
	ass.NoError(err)
	module, err := NewUserModule(db, zap.NewNop(), &util.Identifier{})
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
		fmt.Println("INPUT", input)
		outcome, err := module.Mutation.UserCreate(context.Background(), input)
		ass.NoError(err)
		fmt.Println("OUTCOME: ", outcome)
	})

	// email duplication
}
