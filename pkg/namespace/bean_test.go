package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	log "gorm.io/gorm/logger"

	"bean/pkg/namespace/api/fixtures"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/user"
	uFixtures "bean/pkg/user/api/fixtures"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/connect"
)

func bean() *NamespaceBean {
	config := &struct {
		Beans struct {
			Namespace *Genetic `yaml:"namespace"`
		} `yaml:"beans"`
	}{}

	err := util.ParseFile("../../config.yaml", &config)
	if nil != err {
		panic(err)
	}

	db := util.MockDatabase()
	db.Logger.LogMode(log.Silent)
	logger := util.MockLogger()
	id := util.MockIdentifier()
	bUser := user.NewUserBean(db, logger, id)
	this := NewNamespaceBean(db, logger, id, bUser, config.Beans.Namespace)

	return this
}

func tearDown(bean *NamespaceBean) {
	bean.db.Table(connect.TableNamespaceDomains).Where("id != ?", "").Delete(&model.DomainName{})
	bean.db.Table(connect.TableUserEmail).Where("id != ?", "").Delete(&mUser.UserEmail{})
}

func Test_Namespace(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	util.MockInstall(this, this.db)
	iCreate := fixtures.NamespaceCreateInputFixture(false)

	t.Run("Create", func(t *testing.T) {
		defer tearDown(this)
		t.Run("happy case", func(t *testing.T) {
			now := time.Now()

			claims := &util.Claims{
				StandardClaims: jwt.StandardClaims{
					Audience: this.id.MustULID(),
					Subject:  this.id.MustULID(),
				},
				Kind: util.KindAuthenticated,
			}
			ctx := context.WithValue(context.Background(), util.CxtKeyClaims, claims)
			out, err := this.NamespaceCreate(ctx, iCreate)
			ass.NoError(err)
			ass.Nil(out.Errors)
			ass.Equal(model.NamespaceKindOrganisation, out.Namespace.Kind)
			ass.Equal(*iCreate.Object.Title, out.Namespace.Title)
			ass.Equal(iCreate.Object.IsActive, out.Namespace.IsActive)
			ass.True(out.Namespace.CreatedAt.UnixNano() >= now.UnixNano())
			ass.True(out.Namespace.UpdatedAt.UnixNano() >= now.UnixNano())
			ass.Equal(out.Namespace.Language, api.LanguageAU)

			// check that owner role is created
			// -------
			ownerNS := &model.Namespace{}
			err = this.db.First(&ownerNS, "parent_id = ?", out.Namespace.ID).Error
			ass.NoError(err)
			ass.Equal(ownerNS.Title, "owner")
			ass.Equal(ownerNS.Kind, model.NamespaceKindRole)
			ass.Equal(ownerNS.Language, api.LanguageDefault)

			// check that memberships are setup correctly.
			var counter int64
			this.db.
				Table(connect.TableNamespaceMemberships).
				Where("user_id = ? AND namespace_id = ?", claims.UserId(), out.Namespace.ID).
				Count(&counter)
			ass.Equal(int64(1), counter)

			this.db.
				Table(connect.TableNamespaceMemberships).
				Where("user_id = ? AND namespace_id = ?", claims.UserId(), ownerNS.ID).
				Count(&counter)
			ass.Equal(int64(1), counter)
		})

		t.Run("domain duplication", func(t *testing.T) {
			// create again with same input
			outcome, err := this.NamespaceCreate(context.Background(), iCreate)

			ass.Nil(outcome)
			ass.NotNil(err)
			ass.Contains(err.Error(), "UNIQUE constraint failed: namespace_domains.value")
		})
	})

	t.Run("Query", func(t *testing.T) {
		defer tearDown(this)

		// setup data for query
		oCreate, err := this.NamespaceCreate(context.Background(), iCreate)
		ass.NoError(err)
		id := oCreate.Namespace.ID

		t.Run("load by ID", func(t *testing.T) {
			obj, err := this.Load(context.Background(), id)
			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})

		t.Run("load by domain name -> inactive domain name", func(t *testing.T) {
			domainName := util.Uri(*iCreate.Object.DomainNames.Secondary[1].Value)
			obj, err := this.Namespace(context.Background(), dto.NamespaceFilters{
				Domain: &domainName,
			})

			ass.Error(err)
			ass.Equal(err.Error(), "domain name is not active")
			ass.Nil(obj)
		})

		t.Run("load by domain name -> verified", func(t *testing.T) {
			domainName := util.Uri(*iCreate.Object.DomainNames.Primary.Value)
			obj, err := this.Namespace(context.Background(), dto.NamespaceFilters{Domain: &domainName})

			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})
	})

	t.Run("Update", func(t *testing.T) {
		defer tearDown(this)

		// create namespace so we have something to update
		outcome, err := this.NamespaceCreate(context.Background(), iCreate)
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
				obj, err := this.Load(context.Background(), outcome.Namespace.ID)
				ass.NoError(err)
				ass.Equal(obj.Language, api.LanguageUS)
			}

			features, err := this.Resolvers.Object.Features(context.Background(), outcome.Namespace)
			ass.NoError(err)
			ass.True(features.Register)
		})

		t.Run("version conflict", func(t *testing.T) {
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
	})
}

func Test_Membership(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	util.MockInstall(this, this.db)

	// setup data for query
	// -------
	// create namespace
	iNamespace := fixtures.NamespaceCreateInputFixture(false)
	oNamespace, err := this.NamespaceCreate(context.Background(), iNamespace)
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
			namespace, err := this.Load(context.Background(), oNamespace.Namespace.ID)
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
	})

	t.Run("Update", func(t *testing.T) {
		defer tearDown(this)

		// setup data for query
		// -------
		// create namespace
		iNamespace := fixtures.NamespaceCreateInputFixture(true)
		oNamespace, err := this.NamespaceCreate(context.Background(), iNamespace)
		ass.NoError(err)

		// create user
		iUser := uFixtures.NewUserCreateInputFixture()
		oUser, err := this.user.Resolvers.Mutation.UserCreate(context.Background(), iUser)
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
					obj, err := this.Resolvers.Query.Membership(context.Background(), membership.ID, util.NilString("InvalidVersion"))
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
	})
}
