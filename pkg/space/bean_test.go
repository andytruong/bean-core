package space

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	log "gorm.io/gorm/logger"

	"bean/components/claim"
	"bean/components/conf"
	"bean/components/scalar"
	"bean/pkg/space/api/fixtures"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	"bean/pkg/user"
	uFixtures "bean/pkg/user/api/fixtures"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/connect"
)

func bean() *SpaceBean {
	config := &struct {
		Beans struct {
			Space *Genetic `yaml:"space"`
		} `yaml:"beans"`
	}{}

	err := conf.ParseFile("../../config.yaml", &config)
	if nil != err {
		panic(err)
	}

	db := util.MockDatabase()
	db.Logger.LogMode(log.Silent)
	logger := util.MockLogger()
	id := util.MockIdentifier()
	bUser := user.NewUserBean(db, logger, id)
	this := NewSpaceBean(db, logger, id, bUser, config.Beans.Space)

	return this
}

func tearDown(bean *SpaceBean) {
	bean.db.Model(model.DomainName{}).Where("id != ?", "").Delete(&model.DomainName{})
	bean.db.Table(connect.TableUserEmail).Where("id != ?", "").Delete(&mUser.UserEmail{})
}

func Test_Space(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	util.MockInstall(this, this.db)
	iCreate := fixtures.SpaceCreateInputFixture(false)

	t.Run("Create", func(t *testing.T) {
		defer tearDown(this)
		t.Run("happy case", func(t *testing.T) {
			now := time.Now()

			claims := &claim.Payload{
				StandardClaims: jwt.StandardClaims{
					Audience: this.id.MustULID(),
					Subject:  this.id.MustULID(),
				},
				Kind: claim.KindAuthenticated,
			}
			ctx := context.WithValue(context.Background(), claim.ContextKey, claims)
			out, err := this.SpaceCreate(ctx, iCreate)
			ass.NoError(err)
			ass.Nil(out.Errors)
			ass.Equal(model.SpaceKindOrganisation, out.Space.Kind)
			ass.Equal(*iCreate.Object.Title, out.Space.Title)
			ass.Equal(iCreate.Object.IsActive, out.Space.IsActive)
			ass.True(out.Space.CreatedAt.UnixNano() >= now.UnixNano())
			ass.True(out.Space.UpdatedAt.UnixNano() >= now.UnixNano())
			ass.Equal(out.Space.Language, api.LanguageAU)

			// check that owner role is created
			// -------
			ownerNS := &model.Space{}
			err = this.db.First(&ownerNS, "parent_id = ?", out.Space.ID).Error
			ass.NoError(err)
			ass.Equal(ownerNS.Title, "owner")
			ass.Equal(ownerNS.Kind, model.SpaceKindRole)
			ass.Equal(ownerNS.Language, api.LanguageDefault)

			// check that memberships are setup correctly.
			var counter int64
			this.db.
				Model(&model.Membership{}).
				Where("user_id = ? AND space_id = ?", claims.UserId(), out.Space.ID).
				Count(&counter)
			ass.Equal(int64(1), counter)

			this.db.
				Model(&model.Membership{}).
				Where("user_id = ? AND space_id = ?", claims.UserId(), ownerNS.ID).
				Count(&counter)
			ass.Equal(int64(1), counter)
		})

		t.Run("domain duplication", func(t *testing.T) {
			// create again with same input
			outcome, err := this.SpaceCreate(context.Background(), iCreate)

			ass.Nil(outcome)
			ass.NotNil(err)
			ass.Contains(err.Error(), "UNIQUE constraint failed: space_domains.value")
		})
	})

	t.Run("Query", func(t *testing.T) {
		defer tearDown(this)

		// setup data for query
		oCreate, err := this.SpaceCreate(context.Background(), iCreate)
		ass.NoError(err)
		id := oCreate.Space.ID

		t.Run("load by ID", func(t *testing.T) {
			obj, err := this.Load(context.Background(), id)
			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})

		t.Run("load by domain name -> inactive domain name", func(t *testing.T) {
			domainName := scalar.Uri(*iCreate.Object.DomainNames.Secondary[1].Value)
			obj, err := this.Space(context.Background(), dto.SpaceFilters{
				Domain: &domainName,
			})

			ass.Error(err)
			ass.Equal(err.Error(), "domain name is not active")
			ass.Nil(obj)
		})

		t.Run("load by domain name -> verified", func(t *testing.T) {
			domainName := scalar.Uri(*iCreate.Object.DomainNames.Primary.Value)
			obj, err := this.Space(context.Background(), dto.SpaceFilters{Domain: &domainName})

			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})
	})

	t.Run("Update", func(t *testing.T) {
		defer tearDown(this)

		// create space so we have something to update
		outcome, err := this.SpaceCreate(context.Background(), iCreate)
		ass.NoError(err)
		ass.Nil(outcome.Errors)

		t.Run("happy case", func(t *testing.T) {
			_, err = this.SpaceUpdate(context.Background(), dto.SpaceUpdateInput{
				SpaceID:      outcome.Space.ID,
				SpaceVersion: outcome.Space.Version,
				Object: &dto.SpaceUpdateInputObject{
					Language: api.LanguageUS.Nil(),
					Features: &dto.SpaceUpdateInputFeatures{
						Register: scalar.NilBool(true),
					},
				},
			})

			{
				obj, err := this.Load(context.Background(), outcome.Space.ID)
				ass.NoError(err)
				ass.Equal(obj.Language, api.LanguageUS)
			}

			features, err := this.Resolvers.Object.Features(context.Background(), outcome.Space)
			ass.NoError(err)
			ass.True(features.Register)
		})

		t.Run("version conflict", func(t *testing.T) {
			_, err = this.SpaceUpdate(context.Background(), dto.SpaceUpdateInput{
				SpaceID:      outcome.Space.ID,
				SpaceVersion: "invalid-version",
				Object: &dto.SpaceUpdateInputObject{
					Features: &dto.SpaceUpdateInputFeatures{
						Register: scalar.NilBool(true),
					},
				},
			})

			ass.Equal(err, util.ErrorVersionConflict)
		})
	})
}

func Test_Membership(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	util.MockInstall(this, this.db)

	// setup data for query
	// -------
	// create space
	iSpace := fixtures.SpaceCreateInputFixture(false)
	oSpace, err := this.SpaceCreate(context.Background(), iSpace)
	ass.NoError(err)

	// create user
	iUser := uFixtures.NewUserCreateInputFixture()
	oUser, err := this.user.Resolvers.Mutation.UserCreate(context.Background(), iUser)
	ass.NoError(err)

	t.Run("Create", func(t *testing.T) {
		defer tearDown(this)

		t.Run("create membership", func(t *testing.T) {
			// change feature ON
			{
				oUpdate, err := this.SpaceUpdate(context.Background(), dto.SpaceUpdateInput{
					SpaceID:      oSpace.Space.ID,
					SpaceVersion: oSpace.Space.Version,
					Object: &dto.SpaceUpdateInputObject{
						Features: &dto.SpaceUpdateInputFeatures{
							Register: scalar.NilBool(true),
						},
					},
				})

				ass.NoError(err)
				ass.NotNil(oUpdate)
				ass.NotEqual(oSpace.Space.Version, oUpdate.Space.Version)
			}

			input := dto.SpaceMembershipCreateInput{
				SpaceID:  oSpace.Space.ID,
				UserID:   oUser.User.ID,
				IsActive: false,
			}

			outcome, err := this.SpaceMembershipCreate(context.Background(), input)

			ass.NoError(err)
			ass.Len(outcome.Errors, 0)
			ass.Equal(outcome.Membership.SpaceID, oSpace.Space.ID)
			ass.False(outcome.Membership.IsActive)
		})

		t.Run("create failed of feature is off", func(t *testing.T) {
			space, err := this.Load(context.Background(), oSpace.Space.ID)
			ass.NoError(err)

			// change feature off
			{
				oUpdate, err := this.SpaceUpdate(context.Background(), dto.SpaceUpdateInput{
					SpaceID:      space.ID,
					SpaceVersion: space.Version,
					Object: &dto.SpaceUpdateInputObject{
						Features: &dto.SpaceUpdateInputFeatures{
							Register: scalar.NilBool(false),
						},
					},
				})

				ass.NoError(err)
				ass.NotNil(oUpdate)
				ass.Nil(oUpdate.Errors)
				ass.NotEqual(space.Version, oUpdate.Space.Version)
			}

			// create
			input := dto.SpaceMembershipCreateInput{
				SpaceID:  oSpace.Space.ID,
				UserID:   oUser.User.ID,
				IsActive: false,
			}

			outcome, err := this.SpaceMembershipCreate(
				context.Background(),
				input,
			)

			// check error
			ass.Contains(err.Error(), util.ErrorConfig.Error())
			ass.Contains(err.Error(), "register is off")
			ass.Nil(outcome)
		})
	})

	t.Run("Update", func(t *testing.T) {
		defer tearDown(this)

		// setup data for query
		// -------
		// create space
		iSpace := fixtures.SpaceCreateInputFixture(true)
		oSpace, err := this.SpaceCreate(context.Background(), iSpace)
		ass.NoError(err)

		// create user
		iUser := uFixtures.NewUserCreateInputFixture()
		oUser, err := this.user.Resolvers.Mutation.UserCreate(context.Background(), iUser)
		ass.NoError(err)

		t.Run("create membership", func(t *testing.T) {
			input := dto.SpaceMembershipCreateInput{
				SpaceID:  oSpace.Space.ID,
				UserID:   oUser.User.ID,
				IsActive: false,
			}

			_, err := this.SpaceMembershipCreate(context.Background(), input)
			ass.NoError(err)
		})

		t.Run("update membership", func(t *testing.T) {
			membership := &model.Membership{}

			// create a membership with status OFF.
			{
				input := dto.SpaceMembershipCreateInput{
					SpaceID:  oSpace.Space.ID,
					UserID:   oUser.User.ID,
					IsActive: false,
				}

				outcome, err := this.SpaceMembershipCreate(context.Background(), input)
				ass.NoError(err)
				membership = outcome.Membership
			}

			// load membership
			{
				// without version
				{
					obj, err := this.Resolvers.Query.Membership(context.Background(), membership.ID, nil)
					ass.NoError(err)
					ass.False(obj.IsActive)
				}

				// with version
				{
					obj, err := this.Resolvers.Query.Membership(context.Background(), membership.ID, &membership.Version)
					ass.NoError(err)
					ass.False(obj.IsActive)
				}

				// with invalid version
				{
					obj, err := this.Resolvers.Query.Membership(context.Background(), membership.ID, scalar.NilString("InvalidVersion"))
					ass.Error(err)
					ass.Equal(err.Error(), util.ErrorVersionConflict.Error())
					ass.Nil(obj)
				}
			}

			// change status to ON
			{
				outcome, err := this.SpaceMembershipUpdate(
					context.Background(),
					dto.SpaceMembershipUpdateInput{
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
	})
}
