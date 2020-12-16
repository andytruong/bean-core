package user

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	util2 "bean/components/util"
	connect2 "bean/components/util/connect"
	"bean/pkg/user/api/fixtures"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

func tearDown(this *UserBundle) {
	this.db.Table(connect2.TableUserEmail).Where("1").Delete(&model.UserEmail{})
}

func Test(t *testing.T) {
	ass := assert.New(t)
	db := util2.MockDatabase()
	this := NewUserBundle(db, util2.MockLogger(), util2.MockIdentifier())
	util2.MockInstall(this, db)
	iCreate := fixtures.NewUserCreateInputFixture()

	t.Run("Create", func(t *testing.T) {
		defer tearDown(this)

		t.Run("test happy case, no error", func(t *testing.T) {
			now := time.Now()

			resolver := this.resolvers["UserMutation"].(map[string]interface{})["Create"].(func(context.Context, *dto.UserCreateInput) (*dto.UserMutationOutcome, error))
			out, err := resolver(context.Background(), iCreate)
			ass.NoError(err)
			ass.Empty(out.Errors)
			ass.Equal("https://foo.bar", string(*out.User.AvatarURI))

			{
				resolver := this.resolvers["UserQuery"].(map[string]interface{})["Load"].(func(ctx context.Context, id string) (*model.User, error))
				theUser, err := resolver(context.Background(), out.User.ID)
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
	})
}

func Test_Update(t *testing.T) {
	ass := assert.New(t)
	db := util2.MockDatabase()
	this := NewUserBundle(db, util2.MockLogger(), util2.MockIdentifier())
	util2.MockInstall(this, db)
	iCreate := fixtures.NewUserCreateInputFixture()

	t.Run("Update", func(t *testing.T) {
		defer tearDown(this)

		// create user so we can edit
		rCreate := this.resolvers["UserMutation"].(map[string]interface{})["Create"].(func(context.Context, *dto.UserCreateInput) (*dto.UserMutationOutcome, error))
		oCreate, err := rCreate(context.Background(), iCreate)
		ass.NoError(err)
		ass.NotNil(oCreate)

		t.Run("version conflict", func(t *testing.T) {
			rUpdate := this.resolvers["UserMutation"].(map[string]interface{})["Update"].(func(context.Context, dto.UserUpdateInput) (*dto.UserMutationOutcome, error))
			oUpdate, err := rUpdate(context.Background(), dto.UserUpdateInput{
				ID:      oCreate.User.ID,
				Version: this.id.MustULID(), // some other version
			})

			ass.NoError(err)
			ass.Equal(oUpdate.Errors[0].Code.String(), util2.ErrorCodeConflict.String())
		})

		t.Run("update password", func(t *testing.T) {
			rUpdate := this.resolvers["UserMutation"].(map[string]interface{})["Update"].(func(context.Context, dto.UserUpdateInput) (*dto.UserMutationOutcome, error))
			oUpdate, err := rUpdate(context.Background(), dto.UserUpdateInput{
				ID:      oCreate.User.ID,
				Version: oCreate.User.Version,
				Values: &dto.UserUpdateValuesInput{
					Password: &dto.UserPasswordInput{
						HashedValue: this.id.MustULID(),
					},
				},
			})

			ass.NoError(err)
			ass.NotNil(oUpdate)
			ass.NotEqual(oCreate.User.Version, oUpdate.User.Version)
		})
	})
}
