package user

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"bean/pkg/user/api/fixtures"
	"bean/pkg/user/model/dto"
	"bean/pkg/util"
)

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	this := NewUserBean(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(this, db)
	input := fixtures.NewUserCreateInputFixture()

	t.Run("test happy case, no error", func(t *testing.T) {
		now := time.Now()
		out, err := this.UserCreate(context.Background(), input)
		ass.NoError(err)
		ass.Empty(out.Errors)
		ass.Equal("https://foo.bar", string(*out.User.AvatarURI))

		{
			theUser, err := this.User(context.Background(), out.User.ID)
			ass.NoError(err)
			ass.True(theUser.CreatedAt.UnixNano() >= now.UnixNano())
			ass.Equal("https://foo.bar", string(*theUser.AvatarURI))
			ass.Equal(26, len(theUser.ID))
			ass.Equal(26, len(theUser.Version))
		}
	})

	t.Run("error by email duplication", func(t *testing.T) {
		ass.True(true, "TODO")
	})

	t.Run("email is casted to lower-case", func(t *testing.T) {
		ass.True(true, "TODO")
	})
}

func Test_Update(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	this := NewUserBean(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(this, db)
	input := fixtures.NewUserCreateInputFixture()

	// create user so we can edit
	oCreate, err := this.UserCreate(context.Background(), input)
	ass.NoError(err)
	ass.NotNil(oCreate)

	t.Run("version conflict", func(t *testing.T) {
		oUpdate, err := this.UserUpdate(context.Background(), dto.UserUpdateInput{
			ID:      oCreate.User.ID,
			Version: this.id.MustULID(), // some other version
		})

		ass.NoError(err)
		ass.Equal(oUpdate.Errors[0].Code.String(), util.ErrorCodeConflict.String())
	})

	t.Run("update password", func(t *testing.T) {
		oUpdate, err := this.UserUpdate(context.Background(), dto.UserUpdateInput{
			ID:      oCreate.User.ID,
			Version: oCreate.User.Version,
			Values: &dto.UserUpdateValuesInput{
				Password: &dto.UserPasswordInput{
					Algorithm:   "xxxxxx",
					HashedValue: this.id.MustULID(),
				},
			},
		})

		ass.NoError(err)
		ass.NotNil(oUpdate)
		ass.NotEqual(oCreate.User.Version, oUpdate.User.Version)
	})
}
