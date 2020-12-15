package access

import (
	"context"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"

	"bean/components/claim"
	"bean/components/conf"
	"bean/pkg/access/api/fixtures"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/space"
	fSpace "bean/pkg/space/api/fixtures"
	"bean/pkg/user"
	fUser "bean/pkg/user/api/fixtures"
	"bean/pkg/util"
)

func accessBundle() *AccessBundle {
	config := &struct {
		Bundles struct {
			Access *AccessConfiguration      `yaml:"access"`
			Space  *space.SpaceConfiguration `yaml:"space"`
		} `yaml:"bundles"`
	}{}

	err := conf.ParseFile("../../config.yaml", config)
	if nil != err {
		panic(err)
	}

	db := util.MockDatabase()
	logger := util.MockLogger()
	id := util.MockIdentifier()
	userBundle := user.NewUserBundle(db, logger, id)
	spaceBundle := space.NewSpaceBundle(db, logger, id, userBundle, config.Bundles.Space)
	bundle := NewAccessBundle(db, id, logger, userBundle, spaceBundle, config.Bundles.Access)
	util.MockInstall(bundle, db)

	return bundle
}

func Test_Config(t *testing.T) {
	ass := assert.New(t)
	this := accessBundle()
	key, err := this.config.GetSignKey()
	ass.NoError(err)
	ass.NotNil(key)
}

func Test_Create(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := accessBundle()

	// create userBundle
	iUser := fUser.NewUserCreateInputFixture()
	oUser, err := this.userBundle.Service.Create(this.db, iUser)
	ass.NoError(err)

	// create space
	ctx = context.WithValue(ctx, claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{Subject: oUser.User.ID},
		Kind:           claim.KindAuthenticated,
	})
	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, err := this.spaceBundle.Service.Create(this.db.WithContext(ctx), iSpace)
	ass.NoError(err)

	t.Run("use credentials", func(t *testing.T) {
		t.Run("email inactive", func(t *testing.T) {
			in := fixtures.SessionCreateInputFixtureUseCredentials(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			in.UseCredentials.Email = iUser.Emails.Secondary[1].Value
			_, err := this.SessionCreate(ctx, in)
			ass.Equal(err.Error(), "userBundle not found")
		})

		t.Run("password unmatched", func(t *testing.T) {
			in := fixtures.SessionCreateInputFixtureUseCredentials(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			in.UseCredentials.HashedPassword = "invalid-password"
			outcome, err := this.SessionCreate(ctx, in)

			ass.NoError(err)
			ass.Equal(util.ErrorCodeInput, *outcome.Errors[0].Code)
			ass.Equal(outcome.Errors[0].Message, "invalid password")
			ass.Equal(outcome.Errors[0].Fields, []string{"input.spaceId"})
		})

		t.Run("ok", func(t *testing.T) {
			in := fixtures.SessionCreateInputFixtureUseCredentials(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			out, err := this.SessionCreate(ctx, in)
			ass.NoError(err)
			ass.Equal(oUser.User.ID, out.Session.UserId)
			ass.Equal(oSpace.Space.ID, out.Session.SpaceId)
			ass.Len(out.Errors, 0)

			// check that code challenged & method are saved correctly
			{
				session, err := this.sessionService.LoadByToken(ctx, this.db, *out.Token)
				ass.NoError(err)
				ass.Equal("S256", session.CodeChallengeMethod)
				ass.Equal(
					"84d89877f0d4041efb6bf91a16f0248f2fd573e6af05c19f96bedb9f882f7882",
					session.CodeChallenge,
				)
			}

			// check that with outcome.Session we can generate JWT
			{
				resolver := this.resolvers["Session"].(map[string]interface{})["Jwt"].(func(context.Context, *model.Session, string) (string, error))
				signedString, err := resolver(ctx, out.Session, "0123456789")

				ass.NoError(err)
				ass.Contains(signedString, "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.")

				// check that JWT is valid
				{
					claims, err := this.JwtService.Validate(signedString)
					ass.NoError(err)
					ass.NotNil(claims)
					ass.Equal(claims.SessionId(), out.Session.ID)
					ass.Equal(claims.UserId(), out.Session.UserId)
					ass.Equal(claims.SpaceId(), out.Session.SpaceId)
					ass.Equal(claims.Roles, []string{"owner"})
					ass.Equal(claims.Kind, out.Session.Kind)
				}
			}
		})
	})

	t.Run("OTLT - one time login token", func(t *testing.T) {
		t.Run("generate", func(t *testing.T) {
			oGenerate, err := this.SessionCreate(ctx, &dto.SessionCreateInput{
				GenerateOTLT: &dto.SessionCreateGenerateOTLT{
					SpaceID: oSpace.Space.ID,
					UserID:  oUser.User.ID,
				},
			})

			ass.NoError(err)
			ass.Equal(claim.KindOTLT, oGenerate.Session.Kind)

			{
				out, err := this.SessionCreate(ctx, &dto.SessionCreateInput{
					UseOTLT: &dto.SessionCreateUseOTLT{Token: *oGenerate.Token},
				})

				ass.NoError(err)
				ass.Equal(claim.KindAuthenticated, out.Session.Kind)

				// load again -> should not be found
				{
					otltSession, err := this.Session(ctx, *oGenerate.Token)
					ass.Error(err)
					ass.Nil(otltSession)
				}
			}
		})
	})
}

func Test_SessionCreate_MembershipNotFound(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := accessBundle()

	// create userBundle
	iUser := fUser.NewUserCreateInputFixture()

	// create space
	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, err := this.spaceBundle.Service.Create(this.db, iSpace)
	ass.NoError(err)

	// base input
	in := fixtures.SessionCreateInputFixtureUseCredentials(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	outcome, err := this.SessionCreate(ctx, in)
	ass.Error(err)
	ass.Nil(outcome)
	ass.Contains(err.Error(), "userBundle not found")
}

func Test_Query(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := accessBundle()

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := this.userBundle.Service.Create(this.db, iUser)

	ctx = context.WithValue(ctx, claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{Subject: oUser.User.ID},
		Kind:           claim.KindAuthenticated,
	})
	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, _ := this.spaceBundle.Service.Create(this.db.WithContext(ctx), iSpace)
	in := fixtures.SessionCreateInputFixtureUseCredentials(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	outcome, err := this.SessionCreate(ctx, in)
	ass.NoError(err)

	// can load session without issue
	session, err := this.Session(ctx, *outcome.Token)
	ass.NoError(err)
	ass.Equal(session.SpaceId, oSpace.Space.ID)
	ass.Equal(session.UserId, oUser.User.ID)

	t.Run("load expired session", func(t *testing.T) {
		// change session expiration time
		oneMinDuration, _ := time.ParseDuration("129h")
		session.ExpiredAt = session.ExpiredAt.Add(-1 * oneMinDuration)
		err := this.db.Save(&session).Error
		ass.NoError(err)

		// load again -> error: session expired
		_, err = this.Session(ctx, *outcome.Token)
		ass.Error(err)
		ass.Equal(err.Error(), "session expired")
	})

	t.Run("load one-time-login session -> session deleted", func(t *testing.T) {
		// TODO
	})
}

func Test_Archive(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := accessBundle()

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := this.userBundle.Service.Create(this.db, iUser)
	ctx = context.WithValue(ctx, claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{Subject: oUser.User.ID},
		Kind:           claim.KindAuthenticated,
	})
	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, _ := this.spaceBundle.Service.Create(this.db.WithContext(ctx), iSpace)
	in := fixtures.SessionCreateInputFixtureUseCredentials(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	sessionOutcome, err := this.SessionCreate(ctx, in)
	ass.NoError(err)

	{
		ctx = context.WithValue(context.Background(), claim.ContextKey, &claim.Payload{
			StandardClaims: jwt.StandardClaims{Id: sessionOutcome.Session.ID, Subject: oUser.User.ID},
			Kind:           claim.KindAuthenticated,
		})

		// can archive session without issue
		{
			outcome, err := this.SessionArchive(ctx)
			ass.NoError(err)
			ass.Equal(outcome.Result, true)
		}

		// archive again -> should have error
		{
			outcome, err := this.SessionArchive(ctx)
			ass.NoError(err)
			ass.Equal(outcome.Result, false)
		}
	}
}
