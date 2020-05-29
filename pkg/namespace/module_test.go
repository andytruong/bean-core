package namespace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/pkg/namespace/api/fixtures"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/user"
	uFixtures "bean/pkg/user/api/fixtures"
	"bean/pkg/util"
)

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	input := fixtures.NamespaceCreateInputFixture()
	db := util.MockDatabase().LogMode(false)
	logger := util.MockLogger()
	id := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, id)
	module := NewNamespaceModule(db, logger, id, mUser)
	util.MockInstall(module, db)

	t.Run("happy case", func(t *testing.T) {
		now := time.Now()
		outcome, err := module.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		ass.Nil(outcome.Errors)
		ass.Equal(input.Object.Title, outcome.Namespace.Title)
		ass.Equal(input.Object.IsActive, outcome.Namespace.IsActive)
		ass.True(outcome.Namespace.CreatedAt.UnixNano() >= now.UnixNano())
		ass.True(outcome.Namespace.UpdatedAt.UnixNano() >= now.UnixNano())
	})

	t.Run("domain duplication", func(t *testing.T) {
		// create again with same input
		outcome, err := module.NamespaceCreate(context.Background(), input)

		ass.Nil(outcome)
		ass.NotNil(err)
		ass.Contains(err.Error(), "UNIQUE constraint failed: namespace_domains.value")
	})
}

func Test_Query(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	logger := util.MockLogger()
	identifier := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, identifier)
	module := NewNamespaceModule(db, logger, identifier, mUser)
	util.MockInstall(module, db)

	var id string
	input := fixtures.NamespaceCreateInputFixture()

	{
		// setup data for query
		outcome, err := module.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		id = outcome.Namespace.ID
	}

	{
		obj, err := module.Namespace(context.Background(), id)
		ass.NoError(err)
		ass.Equal(obj.ID, id)
		ass.Equal(obj.Title, input.Object.Title)
		ass.Equal(obj.IsActive, input.Object.IsActive)
	}
}

func Test_Update(t *testing.T) {
	ass := assert.New(t)
	input := fixtures.NamespaceCreateInputFixture()
	db := util.MockDatabase()
	logger := util.MockLogger()
	identifier := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, identifier)
	module := NewNamespaceModule(db, logger, identifier, mUser)
	util.MockInstall(module, db)

	// create namespace so we have something to update
	outcome, err := module.NamespaceCreate(context.Background(), input)
	ass.NoError(err)
	ass.Nil(outcome.Errors)

	t.Run("happy case", func(t *testing.T) {
		_, err = module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
			NamespaceID:      outcome.Namespace.ID,
			NamespaceVersion: outcome.Namespace.Version,
			Object: &dto.NamespaceUpdateInputObject{
				Features: &dto.NamespaceUpdateInputFeatures{
					Register: util.NilBool(true),
				},
			},
		})

		features, err := module.Features(context.Background(), outcome.Namespace)
		ass.NoError(err)
		ass.True(features.Register)
	})

	t.Run("update with invalid version -> conflict", func(t *testing.T) {
		_, err = module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
			NamespaceID:      outcome.Namespace.ID,
			NamespaceVersion: "invalid-version",
			Object: &dto.NamespaceUpdateInputObject{
				Features: &dto.NamespaceUpdateInputFeatures{
					Register: util.NilBool(true),
				},
			},
		})

		ass.Equal(err, util.ErrorVersionConflict)
	})
}

func Test_Membership(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	logger := util.MockLogger()
	identifier := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, identifier)
	module := NewNamespaceModule(db, logger, identifier, mUser)
	util.MockInstall(module, db)

	// setup data for query
	// -------
	// create namespace
	iNamespace := fixtures.NamespaceCreateInputFixture()
	oNamespace, err := module.NamespaceCreate(context.Background(), iNamespace)
	ass.NoError(err)

	// create user
	iUser := uFixtures.NewUserCreateInputFixture()
	oUser, err := mUser.UserCreate(context.Background(), iUser)
	ass.NoError(err)

	t.Run("create membership", func(t *testing.T) {
		// change feature ON
		{
			ok, err := module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
				NamespaceID:      oNamespace.Namespace.ID,
				NamespaceVersion: oNamespace.Namespace.Version,
				Object: &dto.NamespaceUpdateInputObject{
					Features: &dto.NamespaceUpdateInputFeatures{
						Register: util.NilBool(true),
					},
				},
			})

			ass.NoError(err)
			ass.True(*ok)
		}

		input := dto.NamespaceMembershipCreateInput{
			NamespaceID: oNamespace.Namespace.ID,
			UserID:      oUser.User.ID,
			IsActive:    false,
		}

		outcome, err := module.NamespaceMembershipCreate(context.Background(), input)

		ass.NoError(err)
		ass.Len(outcome.Errors, 0)
		ass.Equal(outcome.Membership.NamespaceID, oNamespace.Namespace.ID)
		ass.False(outcome.Membership.IsActive)
	})

	t.Run("create failed of feature is off", func(t *testing.T) {
		namespace, err := module.Namespace(context.Background(), oNamespace.Namespace.ID)
		ass.NoError(err)

		// change feature off
		{
			ok, err := module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
				NamespaceID:      namespace.ID,
				NamespaceVersion: namespace.Version,
				Object: &dto.NamespaceUpdateInputObject{
					Features: &dto.NamespaceUpdateInputFeatures{
						Register: util.NilBool(false),
					},
				},
			})

			ass.NoError(err)
			ass.True(*ok)
		}

		// create
		input := dto.NamespaceMembershipCreateInput{
			NamespaceID: oNamespace.Namespace.ID,
			UserID:      oUser.User.ID,
			IsActive:    false,
		}

		outcome, err := module.NamespaceMembershipCreate(
			context.Background(),
			input,
		)

		// check error
		ass.Contains(err.Error(), util.ErrorConfig.Error())
		ass.Contains(err.Error(), "register is off")
		ass.Nil(outcome)
	})

	t.Run("update membership", func(t *testing.T) {
		membership := &model.Membership{}

		// change feature ON
		{
			ok, err := module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
				NamespaceID:      oNamespace.Namespace.ID,
				NamespaceVersion: oNamespace.Namespace.Version,
				Object: &dto.NamespaceUpdateInputObject{
					Features: &dto.NamespaceUpdateInputFeatures{
						Register: util.NilBool(true),
					},
				},
			})

			ass.NoError(err)
			ass.True(*ok)
		}

		// create a membership with status OFF.
		{
			input := dto.NamespaceMembershipCreateInput{
				NamespaceID: oNamespace.Namespace.ID,
				UserID:      oUser.User.ID,
				IsActive:    false,
			}

			outcome, err := module.NamespaceMembershipCreate(context.Background(), input)
			ass.NoError(err)
			membership = outcome.Membership
		}

		// WIP: change status to ON
		if false {
			fmt.Println("membership: ", membership.ID)
		}
	})
}
