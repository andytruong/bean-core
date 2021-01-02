package space

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	
	"bean/components/claim"
	"bean/components/conf"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
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
	
	lgr := util.MockLogger()
	idr := util.MockIdentifier()
	userBundle := user.NewUserBundle(lgr, idr)
	bundle := NewSpaceBundle(lgr, idr, userBundle, config.Bundles.Space)
	
	return bundle
}

func tearDown(db *gorm.DB) {
	db.Model(model.DomainName{}).Where("id != ?", "").Delete(&model.DomainName{})
	db.Table(connect.TableUserEmail).Where("id != ?", "").Delete(&mUser.UserEmail{})
}

func Test_Space(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)
	iCreate := fixtures.SpaceCreateInputFixture(false)
	
	t.Run("create", func(t *testing.T) {
		defer tearDown(db)
		t.Run("happy case", func(t *testing.T) {
			// setup auth context
			claims := claim.NewPayload()
			claims.
				SetKind(claim.KindAuthenticated).
				SetUserId(bundle.idr.MustULID()).
				SetSpaceId(bundle.idr.MustULID())
			ctx := claim.PayloadToContext(ctx, &claims)
			now := time.Now()
			
			out, err := bundle.Service.Create(ctx, iCreate)
			
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
			err = db.First(&ownerNS, "parent_id = ?", out.Space.ID).Error
			ass.NoError(err)
			ass.Equal(ownerNS.Title, "owner")
			ass.Equal(ownerNS.Kind, model.SpaceKindRole)
			ass.Equal(ownerNS.Language, api.LanguageDefault)
			
			// check that memberships are setup correctly.
			var counter int64
			db.
				Model(&model.Membership{}).
				Where("user_id = ? AND space_id = ?", claims.UserId(), out.Space.ID).
				Count(&counter)
			ass.Equal(int64(1), counter)
			
			db.
				Model(&model.Membership{}).
				Where("user_id = ? AND space_id = ?", claims.UserId(), ownerNS.ID).
				Count(&counter)
			ass.Equal(int64(1), counter)
		})
		
		t.Run("domain duplication", func(t *testing.T) {
			// create again with same input
			outcome, err := bundle.Service.Create(ctx, iCreate)
			
			ass.Nil(outcome)
			ass.NotNil(err)
			ass.Contains(err.Error(), "UNIQUE constraint failed: space_domains.value")
		})
	})
	
	t.Run("query", func(t *testing.T) {
		defer tearDown(db)
		
		// setup data for query
		oCreate, err := bundle.Service.Create(ctx, iCreate)
		ass.NoError(err)
		id := oCreate.Space.ID
		
		t.Run("load by ID", func(t *testing.T) {
			obj, err := bundle.Service.Load(ctx, id)
			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})
		
		t.Run("load by domain name -> inactive domain name", func(t *testing.T) {
			domainName := scalar.Uri(*iCreate.Object.DomainNames.Secondary[1].Value)
			obj, err := bundle.Service.FindOne(ctx, dto.SpaceFilters{Domain: &domainName})
			
			ass.Error(err)
			ass.Equal(err.Error(), "domain name is not active")
			ass.Nil(obj)
		})
		
		t.Run("load by domain name -> verified", func(t *testing.T) {
			domainName := scalar.Uri(*iCreate.Object.DomainNames.Primary.Value)
			obj, err := bundle.Service.FindOne(ctx, dto.SpaceFilters{Domain: &domainName})
			
			ass.NoError(err)
			ass.Equal(obj.ID, id)
			ass.Equal(obj.Title, *iCreate.Object.Title)
			ass.Equal(obj.IsActive, iCreate.Object.IsActive)
		})
	})
	
	t.Run("update", func(t *testing.T) {
		defer tearDown(db)
		
		// create space so we have something to update
		out, err := bundle.Service.Create(ctx, iCreate)
		ass.NoError(err)
		ass.Nil(out.Errors)
		
		t.Run("happy case", func(t *testing.T) {
			_, err = bundle.Service.Update(ctx, *out.Space, dto.SpaceUpdateInput{
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
				obj, err := bundle.Service.Load(ctx, out.Space.ID)
				ass.NoError(err)
				ass.Equal(obj.Language, api.LanguageUS)
			}
			
			features, err := bundle.configService.List(ctx, out.Space)
			ass.NoError(err)
			ass.True(features.Register)
		})
		
		t.Run("version conflict", func(t *testing.T) {
			_, err = bundle.Service.Update(ctx, *out.Space, dto.SpaceUpdateInput{
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
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)
	
	// setup data for query
	// -------
	// create space
	iSpace := fixtures.SpaceCreateInputFixture(false)
	oSpace, err := bundle.Service.Create(ctx, iSpace)
	ass.NoError(err)
	
	// create user
	iUser := uFixtures.NewUserCreateInputFixture()
	oUser, err := bundle.userBundle.Service.Create(ctx, iUser)
	
	ass.NoError(err)
	
	t.Run("create", func(t *testing.T) {
		defer tearDown(db)
		
		t.Run("create membership", func(t *testing.T) {
			// change feature ON
			{
				oUpdate, err := bundle.Service.Update(ctx, *oSpace.Space, dto.SpaceUpdateInput{
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
			
			out, err := bundle.MemberService.Create(ctx, in)
			ass.NoError(err)
			ass.Len(out.Errors, 0)
			ass.Equal(out.Membership.SpaceID, oSpace.Space.ID)
			ass.False(out.Membership.IsActive)
		})
		
		t.Run("create failed of feature is off", func(t *testing.T) {
			space, err := bundle.Service.Load(ctx, oSpace.Space.ID)
			ass.NoError(err)
			
			// change feature off
			{
				oUpdate, err := bundle.Service.Update(ctx, *space, dto.SpaceUpdateInput{
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
			
			resolver := bundle.resolvers["SpaceMembershipMutation"].(map[string]interface{})["Create"].(func(context.Context, *dto.SpaceMembershipMutation, dto.SpaceMembershipCreateInput) (*dto.SpaceMembershipCreateOutcome, error))
			outcome, err := resolver(ctx, nil, input)
			
			// check error
			ass.Contains(err.Error(), util.ErrorConfig.Error())
			ass.Contains(err.Error(), "register is off")
			ass.Nil(outcome)
		})
	})
	
	t.Run("update", func(t *testing.T) {
		defer tearDown(db)
		
		// setup data for query
		// -------
		// create space
		iSpace := fixtures.SpaceCreateInputFixture(true)
		oSpace, err := bundle.Service.Create(ctx, iSpace)
		ass.NoError(err)
		
		// create user
		iUser := uFixtures.NewUserCreateInputFixture()
		oUser, err := bundle.userBundle.Service.Create(ctx, iUser)
		
		ass.NoError(err)
		
		t.Run("create membership", func(t *testing.T) {
			in := dto.SpaceMembershipCreateInput{
				SpaceID:  oSpace.Space.ID,
				UserID:   oUser.User.ID,
				IsActive: false,
			}
			
			_, err := bundle.MemberService.Create(ctx, in)
			ass.NoError(err)
		})
		
		t.Run("update membership", func(t *testing.T) {
			resolver := bundle.resolvers["SpaceMembershipMutation"].(map[string]interface{})["Update"].(func(context.Context, *dto.SpaceMembershipMutation, dto.SpaceMembershipUpdateInput) (*dto.SpaceMembershipCreateOutcome, error))
			var membership *model.Membership
			
			// create a membership with status OFF.
			{
				in := dto.SpaceMembershipCreateInput{
					SpaceID:  oSpace.Space.ID,
					UserID:   oUser.User.ID,
					IsActive: false,
				}
				
				out, err := bundle.MemberService.Create(ctx, in)
				ass.NoError(err)
				membership = out.Membership
			}
			
			// load membership
			{
				resolver := bundle.resolvers["SpaceMembershipQuery"].(map[string]interface{})["Load"].(func(context.Context, *dto.SpaceMembershipQuery, string, *string) (*model.Membership, error))
				
				// without version
				{
					obj, err := resolver(ctx, nil, membership.ID, nil)
					ass.NoError(err)
					ass.False(obj.IsActive)
				}
				
				// with version
				{
					obj, err := resolver(ctx, nil, membership.ID, &membership.Version)
					ass.NoError(err)
					ass.False(obj.IsActive)
				}
				
				// with invalid version
				{
					obj, err := resolver(ctx, nil, membership.ID, scalar.NilString("InvalidVersion"))
					ass.Error(err)
					ass.Equal(err.Error(), util.ErrorVersionConflict.Error())
					ass.Nil(obj)
				}
			}
			
			// change status to ON
			{
				outcome, err := resolver(
					ctx,
					nil,
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
