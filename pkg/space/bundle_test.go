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
	"bean/components/util"
	"bean/components/util/connect"
	"bean/pkg/infra/api"
	"bean/pkg/space/api/fixtures"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	"bean/pkg/user"
	uFixtures "bean/pkg/user/api/fixtures"
	mUser "bean/pkg/user/model"
)

func bundle() *SpaceBundle {
	config := &struct {
		Bundles struct {
			Space *SpaceConfiguration `yaml:"space"`
		} `yaml:"bundles"`
	}{}

	err := conf.ParseFile("../../config.yaml", &config)
	if nil != err {
		panic(err)
	}

	db := util.MockDatabase()
	db.Logger.LogMode(log.Silent)
	logger := util.MockLogger()
	id := util.MockIdentifier()
	userBundle := user.NewUserBundle(db, logger, id)
	this := NewSpaceBundle(db, logger, id, userBundle, config.Bundles.Space)

	return this
}

func tearDown(bundle *SpaceBundle) {
	bundle.db.Model(model.DomainName{}).Where("id != ?", "").Delete(&model.DomainName{})
	bundle.db.Table(connect.TableUserEmail).Where("id != ?", "").Delete(&mUser.UserEmail{})
}

func Test_Space(t *testing.T) {
	ass := assert.New(t)
	this := bundle()
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
			out, err := this.Service.Create(this.db.WithContext(ctx), iCreate)

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
			outcome, err := this.Service.Create(this.db, iCreate)

			ass.Nil(outcome)
			ass.NotNil(err)
			ass.Contains(err.Error(), "UNIQUE constraint failed: space_domains.value")
		})
	})

	t.Run("Query", func(t *testing.T) {
		defer tearDown(this)

		// setup data for query
		oCreate, err := this.Service.Create(this.db, iCreate)
		ass.NoError(err)
		id := oCreate.Space.ID

		t.Run("load by ID", func(t *testing.T) {
			obj, err := this.Service.Load(context.Background(), id)
			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})

		t.Run("load by domain name -> inactive domain name", func(t *testing.T) {
			domainName := scalar.Uri(*iCreate.Object.DomainNames.Secondary[1].Value)
			obj, err := this.Service.FindOne(context.Background(), dto.SpaceFilters{Domain: &domainName})

			ass.Error(err)
			ass.Equal(err.Error(), "domain name is not active")
			ass.Nil(obj)
		})

		t.Run("load by domain name -> verified", func(t *testing.T) {
			domainName := scalar.Uri(*iCreate.Object.DomainNames.Primary.Value)
			obj, err := this.Service.FindOne(context.Background(), dto.SpaceFilters{Domain: &domainName})

			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})
	})

	t.Run("Update", func(t *testing.T) {
		defer tearDown(this)

		// create space so we have something to update
		out, err := this.Service.Create(this.db, iCreate)
		ass.NoError(err)
		ass.Nil(out.Errors)

		t.Run("happy case", func(t *testing.T) {
			_, err = this.Service.Update(this.db, *out.Space, dto.SpaceUpdateInput{
				SpaceID:      out.Space.ID,
				SpaceVersion: out.Space.Version,
				Object: &dto.SpaceUpdateInputObject{
					Language: api.LanguageUS.Nil(),
					Features: &dto.SpaceUpdateInputFeatures{
						Register: scalar.NilBool(true),
					},
				},
			})

			{
				obj, err := this.Service.Load(context.Background(), out.Space.ID)
				ass.NoError(err)
				ass.Equal(obj.Language, api.LanguageUS)
			}

			features, err := this.configService.List(context.Background(), out.Space)
			ass.NoError(err)
			ass.True(features.Register)
		})

		t.Run("version conflict", func(t *testing.T) {
			_, err = this.Service.Update(this.db, *out.Space, dto.SpaceUpdateInput{
				SpaceID:      out.Space.ID,
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
	this := bundle()
	util.MockInstall(this, this.db)

	// setup data for query
	// -------
	// create space
	iSpace := fixtures.SpaceCreateInputFixture(false)
	oSpace, err := this.Service.Create(this.db, iSpace)
	ass.NoError(err)

	// create user
	iUser := uFixtures.NewUserCreateInputFixture()
	oUser, err := this.userBundle.Service.Create(this.db, iUser)

	ass.NoError(err)

	t.Run("Create", func(t *testing.T) {
		defer tearDown(this)

		t.Run("create membership", func(t *testing.T) {
			// change feature ON
			{
				oUpdate, err := this.Service.Update(this.db, *oSpace.Space, dto.SpaceUpdateInput{
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

			in := dto.SpaceMembershipCreateInput{
				SpaceID:  oSpace.Space.ID,
				UserID:   oUser.User.ID,
				IsActive: false,
			}

			out, err := this.MemberService.Create(this.db, in)
			ass.NoError(err)
			ass.Len(out.Errors, 0)
			ass.Equal(out.Membership.SpaceID, oSpace.Space.ID)
			ass.False(out.Membership.IsActive)
		})

		t.Run("create failed of feature is off", func(t *testing.T) {
			space, err := this.Service.Load(context.Background(), oSpace.Space.ID)
			ass.NoError(err)

			// change feature off
			{
				oUpdate, err := this.Service.Update(this.db, *space, dto.SpaceUpdateInput{
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

			resolver := this.resolvers["SpaceMembershipMutation"].(map[string]interface{})["Create"].(func(context.Context, dto.SpaceMembershipCreateInput) (*dto.SpaceMembershipCreateOutcome, error))
			outcome, err := resolver(context.Background(), input)

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
		oSpace, err := this.Service.Create(this.db, iSpace)
		ass.NoError(err)

		// create user
		iUser := uFixtures.NewUserCreateInputFixture()
		oUser, err := this.userBundle.Service.Create(this.db, iUser)

		ass.NoError(err)

		t.Run("create membership", func(t *testing.T) {
			in := dto.SpaceMembershipCreateInput{
				SpaceID:  oSpace.Space.ID,
				UserID:   oUser.User.ID,
				IsActive: false,
			}

			_, err := this.MemberService.Create(this.db, in)
			ass.NoError(err)
		})

		t.Run("update membership", func(t *testing.T) {
			resolver := this.resolvers["SpaceMembershipMutation"].(map[string]interface{})["Update"].(func(context.Context, dto.SpaceMembershipUpdateInput) (*dto.SpaceMembershipCreateOutcome, error))
			membership := &model.Membership{}

			// create a membership with status OFF.
			{
				in := dto.SpaceMembershipCreateInput{
					SpaceID:  oSpace.Space.ID,
					UserID:   oUser.User.ID,
					IsActive: false,
				}

				out, err := this.MemberService.Create(this.db, in)
				ass.NoError(err)
				membership = out.Membership
			}

			// load membership
			{
				resolver := this.resolvers["SpaceMembershipQuery"].(map[string]interface{})["Load"].(func(context.Context, string, *string) (*model.Membership, error))

				// without version
				{
					obj, err := resolver(context.Background(), membership.ID, nil)
					ass.NoError(err)
					ass.False(obj.IsActive)
				}

				// with version
				{
					obj, err := resolver(context.Background(), membership.ID, &membership.Version)
					ass.NoError(err)
					ass.False(obj.IsActive)
				}

				// with invalid version
				{
					obj, err := resolver(context.Background(), membership.ID, scalar.NilString("InvalidVersion"))
					ass.Error(err)
					ass.Equal(err.Error(), util.ErrorVersionConflict.Error())
					ass.Nil(obj)
				}
			}

			// change status to ON
			{
				outcome, err := resolver(
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
