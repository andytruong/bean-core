package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"bean/pkg/namespace/api/fixtures"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/user"
	uFixtures "bean/pkg/user/api/fixtures"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/connect"
)

func module() *NamespaceModule {
	configRaw, err := util.ParseFile("../../config.yaml")
	if nil != err {
		panic(err)
	}

	config := &struct {
		Modules struct {
			Namespace *Config `yaml:"namespace"`
		} `yaml:"modules"`
	}{}

	{
		err := yaml.Unmarshal(configRaw, &config)
		if nil != err {
			panic(err)
		}
	}

	db := util.MockDatabase().LogMode(false)
	logger := util.MockLogger()
	id := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, id)
	this := NewNamespaceModule(db, logger, id, mUser, config.Modules.Namespace)

	return this
}

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	this := module()
	util.MockInstall(this, this.db)

	input := fixtures.NamespaceCreateInputFixture(false)

	t.Run("happy case", func(t *testing.T) {
		now := time.Now()
		outcome, err := this.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		ass.Nil(outcome.Errors)
		ass.Equal(model.NamespaceKindOrganisation, outcome.Namespace.Kind)
		ass.Equal(*input.Object.Title, outcome.Namespace.Title)
		ass.Equal(input.Object.IsActive, outcome.Namespace.IsActive)
		ass.True(outcome.Namespace.CreatedAt.UnixNano() >= now.UnixNano())
		ass.True(outcome.Namespace.UpdatedAt.UnixNano() >= now.UnixNano())
		ass.Equal(outcome.Namespace.Language, api.LanguageAU)

		// check that owner role is created
		// -------
		ownerNS := &model.Namespace{}
		err = this.db.First(&ownerNS, "parent_id = ?", outcome.Namespace.ID).Error
		ass.NoError(err)
		ass.Equal(ownerNS.Title, "owner")
		ass.Equal(ownerNS.Kind, model.NamespaceKindRole)
		ass.Equal(ownerNS.Language, api.LanguageDefault)

		// check that memberships are setup correctly.
		counter := 0
		this.db.
			Table(connect.TableNamespaceMemberships).
			Where("user_id = ? AND namespace_id = ?", input.Context.UserID, outcome.Namespace.ID).
			Count(&counter)
		ass.Equal(1, counter)

		this.db.
			Table(connect.TableNamespaceMemberships).
			Where("user_id = ? AND namespace_id = ?", input.Context.UserID, ownerNS.ID).
			Count(&counter)
		ass.Equal(1, counter)
	})

	t.Run("domain duplication", func(t *testing.T) {
		// create again with same input
		outcome, err := this.NamespaceCreate(context.Background(), input)

		ass.Nil(outcome)
		ass.NotNil(err)
		ass.Contains(err.Error(), "UNIQUE constraint failed: namespace_domains.value")
	})
}

func Test_Query(t *testing.T) {
	ass := assert.New(t)
	this := module()
	util.MockInstall(this, this.db)

	var id string
	input := fixtures.NamespaceCreateInputFixture(false)

	{
		// setup data for query
		outcome, err := this.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		id = outcome.Namespace.ID
	}

	{
		obj, err := this.Namespace(context.Background(), id)
		ass.NoError(err)
		ass.Equal(obj.ID, id)
		ass.Equal(obj.Title, *input.Object.Title)
		ass.Equal(obj.IsActive, input.Object.IsActive)
	}
}

func Test_Update(t *testing.T) {
	ass := assert.New(t)
	this := module()
	util.MockInstall(this, this.db)

	// create namespace so we have something to update
	input := fixtures.NamespaceCreateInputFixture(false)
	outcome, err := this.NamespaceCreate(context.Background(), input)
	ass.NoError(err)
	ass.Nil(outcome.Errors)

	t.Run("happy case", func(t *testing.T) {
		_, err = this.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
			NamespaceID:      outcome.Namespace.ID,
			NamespaceVersion: outcome.Namespace.Version,
			Object: &dto.NamespaceUpdateInputObject{
				Language: api.LanguageUS.Nil(),
				Features: &dto.NamespaceUpdateInputFeatures{
					Register: util.NilBool(true),
				},
			},
		})

		{
			obj, err := this.Namespace(context.Background(), outcome.Namespace.ID)
			ass.NoError(err)
			ass.Equal(obj.Language, api.LanguageUS)
		}

		features, err := this.Features(context.Background(), outcome.Namespace)
		ass.NoError(err)
		ass.True(features.Register)
	})

	t.Run("update with invalid version -> conflict", func(t *testing.T) {
		_, err = this.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
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

func Test_Membership_Create(t *testing.T) {
	ass := assert.New(t)
	this := module()
	util.MockInstall(this, this.db)

	// setup data for query
	// -------
	// create namespace
	iNamespace := fixtures.NamespaceCreateInputFixture(false)
	oNamespace, err := this.NamespaceCreate(context.Background(), iNamespace)
	ass.NoError(err)

	// create user
	iUser := uFixtures.NewUserCreateInputFixture()
	oUser, err := this.user.UserCreate(context.Background(), iUser)
	ass.NoError(err)

	t.Run("create membership", func(t *testing.T) {
		// change feature ON
		{
			ok, err := this.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
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

		outcome, err := this.NamespaceMembershipCreate(context.Background(), input)

		ass.NoError(err)
		ass.Len(outcome.Errors, 0)
		ass.Equal(outcome.Membership.NamespaceID, oNamespace.Namespace.ID)
		ass.False(outcome.Membership.IsActive)
	})

	t.Run("create failed of feature is off", func(t *testing.T) {
		namespace, err := this.Namespace(context.Background(), oNamespace.Namespace.ID)
		ass.NoError(err)

		// change feature off
		{
			ok, err := this.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
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

		outcome, err := this.NamespaceMembershipCreate(
			context.Background(),
			input,
		)

		// check error
		ass.Contains(err.Error(), util.ErrorConfig.Error())
		ass.Contains(err.Error(), "register is off")
		ass.Nil(outcome)
	})
}

func Test_Membership_Update(t *testing.T) {
	ass := assert.New(t)
	this := module()
	util.MockInstall(this, this.db)

	// setup data for query
	// -------
	// create namespace
	iNamespace := fixtures.NamespaceCreateInputFixture(true)
	oNamespace, err := this.NamespaceCreate(context.Background(), iNamespace)
	ass.NoError(err)

	// create user
	iUser := uFixtures.NewUserCreateInputFixture()
	oUser, err := this.user.UserCreate(context.Background(), iUser)
	ass.NoError(err)

	t.Run("create membership", func(t *testing.T) {
		input := dto.NamespaceMembershipCreateInput{
			NamespaceID: oNamespace.Namespace.ID,
			UserID:      oUser.User.ID,
			IsActive:    false,
		}

		_, err := this.NamespaceMembershipCreate(context.Background(), input)
		ass.NoError(err)
	})

	t.Run("update membership", func(t *testing.T) {
		membership := &model.Membership{}

		// create a membership with status OFF.
		{
			input := dto.NamespaceMembershipCreateInput{
				NamespaceID: oNamespace.Namespace.ID,
				UserID:      oUser.User.ID,
				IsActive:    false,
			}

			outcome, err := this.NamespaceMembershipCreate(context.Background(), input)
			ass.NoError(err)
			membership = outcome.Membership
		}

		// load membership
		{
			// without version
			{
				obj, err := this.Membership(context.Background(), membership.ID, nil)
				ass.NoError(err)
				ass.False(obj.IsActive)
			}

			// with version
			{
				obj, err := this.Membership(context.Background(), membership.ID, &membership.Version)
				ass.NoError(err)
				ass.False(obj.IsActive)
			}

			// with invalid version
			{
				obj, err := this.Membership(context.Background(), membership.ID, util.NilString("InvalidVersion"))
				ass.Error(err)
				ass.Equal(err.Error(), util.ErrorVersionConflict.Error())
				ass.Nil(obj)
			}
		}

		// change status to ON
		{
			outcome, err := this.NamespaceMembershipUpdate(
				context.Background(),
				dto.NamespaceMembershipUpdateInput{
					Id:       membership.ID,
					Version:  membership.Version,
					IsActive: true,
				},
			)

			ass.NoError(err)
			ass.Len(outcome.Errors, 0)
			ass.True(outcome.Membership.IsActive)
			ass.NotEqual(outcome.Membership.Version, membership.Version)
		}
	})
}
