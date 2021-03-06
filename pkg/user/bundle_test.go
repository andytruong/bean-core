package user

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"bean/components/connect"
	"bean/components/util"
	"bean/pkg/user/api/fixtures"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

func tearDown(db *gorm.DB) {
	db.Table(connect.TableUserEmail).Where("1").Delete(&model.UserEmail{})
}

func Test(t *testing.T) {
	ass := assert.New(t)
	db := connect.MockDatabase()
	bundle := NewUserBundle(util.MockLogger(), util.MockIdentifier())
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)
	iCreate := fixtures.NewUserCreateInputFixture()

	t.Run("database schema", func(t *testing.T) {
		ass.True(
			db.Migrator().HasTable("users"),
			"Table `users` was created",
		)
	})

	t.Run("Create", func(t *testing.T) {
		defer tearDown(db)

		t.Run("ok", func(t *testing.T) {
			now := time.Now()

			resolver := bundle.resolvers["UserMutation"].(map[string]interface{})["Create"].(func(context.Context, *dto.UserMutation, *dto.UserCreateInput) (
				*dto.UserMutationOutcome, error,
			))
			out, err := resolver(ctx, nil, iCreate)
			ass.NoError(err)
			ass.Empty(out.Errors)
			ass.Equal("https://foo.bar", string(*out.User.AvatarURI))

			{
				resolver := bundle.resolvers["UserQuery"].(map[string]interface{})["Load"].(func(context.Context, *dto.UserQuery, string) (
					*model.User, error,
				))
				theUser, err := resolver(ctx, nil, out.User.ID)
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
	db := connect.MockDatabase()
	bundle := NewUserBundle(util.MockLogger(), util.MockIdentifier())
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)
	iCreate := fixtures.NewUserCreateInputFixture()

	t.Run("Update", func(t *testing.T) {
		defer tearDown(db)

		// create user so we can edit
		rCreate := bundle.resolvers["UserMutation"].(map[string]interface{})["Create"].(func(context.Context, *dto.UserMutation, *dto.UserCreateInput) (
			*dto.UserMutationOutcome, error,
		))
		oCreate, err := rCreate(ctx, nil, iCreate)
		ass.NoError(err)
		ass.NotNil(oCreate)
		rUpdate := bundle.resolvers["UserMutation"].(map[string]interface{})["Update"].(func(context.Context, *dto.UserMutation, dto.UserUpdateInput) (
			*dto.UserMutationOutcome, error,
		))

		t.Run("version conflict", func(t *testing.T) {
			oUpdate, err := rUpdate(ctx, nil, dto.UserUpdateInput{
				ID:      oCreate.User.ID,
				Version: bundle.idr.ULID(), // some other version
			})

			ass.NoError(err)
			ass.Equal(oUpdate.Errors[0].Code.String(), util.ErrorCodeConflict.String())
		})

		t.Run("update password", func(t *testing.T) {
			oUpdate, err := rUpdate(ctx, nil, dto.UserUpdateInput{
				ID:      oCreate.User.ID,
				Version: oCreate.User.Version,
				Values: &dto.UserUpdateValuesInput{
					Password: &dto.UserPasswordInput{
						HashedValue: bundle.idr.ULID(),
					},
				},
			})

			ass.NoError(err)
			ass.NotNil(oUpdate)
			ass.NotEqual(oCreate.User.Version, oUpdate.User.Version)
		})
	})
}
