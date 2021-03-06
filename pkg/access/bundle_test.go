package access

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/components/claim"
	"bean/components/conf"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/space"
	fSpace "bean/pkg/space/api/fixtures"
	"bean/pkg/user"
	fUser "bean/pkg/user/api/fixtures"
)

func bundle() *Bundle {
	config := &struct {
		Bundles struct {
			Access *Config       `yaml:"access"`
			Space  *space.Config `yaml:"space"`
		} `yaml:"bundles"`
	}{}

	err := conf.ParseFile("../../config.yaml", config)
	if nil != err {
		panic(err)
	}

	lgr := util.MockLogger()
	id := util.MockIdentifier()
	userBundle := user.NewUserBundle(lgr, id)
	spaceBundle := space.NewSpaceBundle(lgr, id, userBundle, config.Bundles.Space)
	bundle, err := NewAccessBundle(id, lgr, userBundle, spaceBundle, config.Bundles.Access)
	if nil != err {
		panic(err)
	}

	return bundle
}

func newCreateSessionInput(spaceId string, email string, hashedPassword string) *dto.SessionCreateInput {
	codeVerifier := []byte("0123456789")

	return &dto.SessionCreateInput{
		SpaceID:             spaceId,
		Email:               scalar.EmailAddress(email),
		HashedPassword:      hashedPassword,
		CodeChallengeMethod: "S256",
		CodeChallenge:       fmt.Sprintf("%x", sha256.Sum256(codeVerifier)),
	}
}

func Test_Config(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	ass.NotNil(bundle.JwtService.privateKey)
}

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	// create userBundle
	iUser := fUser.NewUserCreateInputFixture()
	oUser, err := bundle.userBundle.UserService.Create(ctx, iUser)
	ass.NoError(err)

	// create space
	claims := claim.NewPayload()
	ctx = claims.
		SetUserId(oUser.User.ID).
		SetKind(claim.KindAuthenticated).
		ToContext(ctx)

	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, err := bundle.spaceBundle.Service.Create(ctx, iSpace)
	ass.NoError(err)

	t.Run("use credentials", func(t *testing.T) {
		t.Run("email inactive", func(t *testing.T) {
			in := newCreateSessionInput(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			in.Email = iUser.Emails.Secondary[1].Value
			_, err := bundle.sessionService.newSessionWithCredentials(ctx, in)
			ass.Equal(err, user.ErrorUserNotFound)
		})

		t.Run("password unmatched", func(t *testing.T) {
			in := newCreateSessionInput(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			in.HashedPassword = "invalid-password"
			outcome, err := bundle.sessionService.newSessionWithCredentials(ctx, in)

			ass.NoError(err)
			ass.Equal(util.ErrorCodeInput, *outcome.Errors[0].Code)
			ass.Equal(outcome.Errors[0].Message, "invalid password")
			ass.Equal(outcome.Errors[0].Fields, []string{"input.spaceId"})
		})

		t.Run("ok", func(t *testing.T) {
			in := newCreateSessionInput(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			out, err := bundle.sessionService.newSessionWithCredentials(ctx, in)
			ass.NoError(err)
			ass.Equal(oUser.User.ID, out.Session.UserId)
			ass.Equal(oSpace.Space.ID, out.Session.SpaceId)
			ass.Len(out.Errors, 0)

			// check that code challenged & method are saved correctly
			{
				session, err := bundle.sessionService.LoadByToken(ctx, *out.Token)
				ass.NoError(err)
				ass.Equal("S256", session.CodeChallengeMethod)
				ass.Equal(
					"84d89877f0d4041efb6bf91a16f0248f2fd573e6af05c19f96bedb9f882f7882",
					session.CodeChallenge,
				)
			}

			// check that with outcome.Session we can generate JWT
			{
				resolver := bundle.resolvers["Session"].(map[string]interface{})["Jwt"].(func(context.Context, *model.Session, string) (
					string, error,
				))
				signedString, err := resolver(ctx, out.Session, "0123456789")

				ass.NoError(err)
				ass.Contains(signedString, "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.")

				// check that JWT is valid
				{
					claims, err := bundle.JwtService.Validate(signedString)
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
			oGenerate, err := bundle.sessionService.newOTLTSession(ctx, &dto.SessionCreateOTLTSessionInput{
				SpaceID: oSpace.Space.ID,
				UserID:  oUser.User.ID,
			})

			ass.NoError(err)
			ass.Equal(claim.KindOTLT, oGenerate.Session.Kind)

			{
				out, err := bundle.sessionService.newSessionWithOTLT(ctx, &dto.SessionExchangeOTLTInput{Token: *oGenerate.Token})
				ass.NoError(err)
				ass.Equal(claim.KindAuthenticated, out.Session.Kind)

				// load again -> should not be found
				{
					otltSession, err := bundle.sessionService.LoadByToken(ctx, *oGenerate.Token)
					ass.Error(err)
					ass.Nil(otltSession)
				}
			}
		})
	})
}

func Test_SessionCreate_MembershipNotFound(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	// create userBundle
	iUser := fUser.NewUserCreateInputFixture()

	// create space
	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, err := bundle.spaceBundle.Service.Create(ctx, iSpace)
	ass.NoError(err)

	// base input
	in := newCreateSessionInput(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	out, err := bundle.sessionService.newSessionWithCredentials(ctx, in)
	ass.Error(err)
	ass.Nil(out)
	ass.Contains(err, user.ErrorUserNotFound)
}

func Test_Query(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := bundle.userBundle.UserService.Create(ctx, iUser)

	claims := claim.NewPayload()
	ctx = claims.
		SetUserId(oUser.User.ID).
		SetKind(claim.KindAuthenticated).
		ToContext(ctx)

	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, _ := bundle.spaceBundle.Service.Create(ctx, iSpace)
	in := newCreateSessionInput(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	out, err := bundle.sessionService.newSessionWithCredentials(ctx, in)
	ass.NoError(err)

	// can load session without issue
	session, err := bundle.sessionService.LoadByToken(ctx, *out.Token)
	ass.NoError(err)
	ass.Equal(session.SpaceId, oSpace.Space.ID)
	ass.Equal(session.UserId, oUser.User.ID)

	t.Run("load expired session", func(t *testing.T) {
		// change session expiration time
		oneMinDuration, _ := time.ParseDuration("129h")
		session.ExpiredAt = session.ExpiredAt.Add(-1 * oneMinDuration)
		err := db.Save(&session).Error
		ass.NoError(err)

		// load again -> error: session expired
		_, err = bundle.sessionService.LoadByToken(ctx, *out.Token)
		ass.Error(err)
		ass.Equal(err.Error(), "session expired")
	})

	t.Run("load one-time-login session -> session deleted", func(t *testing.T) {
		// TODO
	})
}

func Test_Archive(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := bundle.userBundle.UserService.Create(ctx, iUser)
	ctx = connect.DBToContext(ctx, db)

	claims := claim.NewPayload()
	ctx = claims.
		SetUserId(oUser.User.ID).
		SetKind(claim.KindAuthenticated).
		ToContext(ctx)

	iSpace := fSpace.SpaceCreateInputFixture(false)
	oSpace, _ := bundle.spaceBundle.Service.Create(ctx, iSpace)
	in := newCreateSessionInput(oSpace.Space.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	sessionOutcome, err := bundle.sessionService.newSessionWithCredentials(ctx, in)
	ass.NoError(err)

	{
		resolver := bundle.resolvers["AccessSessionMutation"].(map[string]interface{})["Archive"].(func(context.Context, *dto.AccessSessionMutation) (
			*dto.SessionArchiveOutcome, error,
		))

		// setup auth context
		claims := claim.NewPayload()
		ctx = claims.
			SetUserId(oUser.User.ID).
			SetKind(claim.KindAuthenticated).
			SetSessionId(sessionOutcome.Session.ID).
			ToContext(ctx)

		// can archive session without issue
		{
			out, err := resolver(ctx, nil)
			ass.NoError(err)
			ass.Equal(out.Result, true)
		}

		// archive again -> should have error
		{
			out, err := resolver(ctx, nil)
			ass.NoError(err)
			ass.Equal(out.Result, false)
		}
	}
}
