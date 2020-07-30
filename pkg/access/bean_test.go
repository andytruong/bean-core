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
	"bean/pkg/access/model/dto"
	"bean/pkg/namespace"
	fNamespace "bean/pkg/namespace/api/fixtures"
	"bean/pkg/user"
	fUser "bean/pkg/user/api/fixtures"
	"bean/pkg/util"
)

func bean() *AccessBean {
	config := &struct {
		Beans struct {
			Access *Genetic `yaml:"access"`
		} `yaml:"beans"`
	}{}

	err := conf.ParseFile("../../config.yaml", config)
	if nil != err {
		panic(err)
	}

	db := util.MockDatabase()
	logger := util.MockLogger()
	id := util.MockIdentifier()
	bUser := user.NewUserBean(db, logger, id)
	bNamespace := namespace.NewNamespaceBean(db, logger, id, bUser, nil)
	bean := NewAccessBean(db, id, logger, bUser, bNamespace, config.Beans.Access)
	util.MockInstall(bean, db)

	return bean
}

func Test_Config(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	key, err := this.genetic.GetSignKey()
	ass.NoError(err)
	ass.NotNil(key)
}

func Test_Create(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := bean()

	// create user
	iUser := fUser.NewUserCreateInputFixture()
	oUser, err := this.user.Resolvers.Mutation.UserCreate(ctx, iUser)
	ass.NoError(err)

	// create namespace
	ctx = context.WithValue(ctx, claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{Subject: oUser.User.ID},
		Kind:           claim.KindAuthenticated,
	})
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	oNamespace, err := this.namespace.NamespaceCreate(ctx, iNamespace)
	ass.NoError(err)

	t.Run("use credentials", func(t *testing.T) {
		t.Run("email inactive", func(t *testing.T) {
			in := fixtures.SessionCreateInputFixtureUseCredentials(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			in.UseCredentials.Email = iUser.Emails.Secondary[1].Value
			_, err := this.SessionCreate(ctx, in)
			ass.Equal(err.Error(), "user not found")
		})

		t.Run("password unmatched", func(t *testing.T) {
			in := fixtures.SessionCreateInputFixtureUseCredentials(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			in.UseCredentials.HashedPassword = "invalid-password"
			outcome, err := this.SessionCreate(ctx, in)

			ass.NoError(err)
			ass.Equal(util.ErrorCodeInput, *outcome.Errors[0].Code)
			ass.Equal(outcome.Errors[0].Message, "invalid password")
			ass.Equal(outcome.Errors[0].Fields, []string{"input.namespaceId"})
		})

		t.Run("ok", func(t *testing.T) {
			in := fixtures.SessionCreateInputFixtureUseCredentials(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
			out, err := this.SessionCreate(ctx, in)
			ass.NoError(err)
			ass.Equal(oUser.User.ID, out.Session.UserId)
			ass.Equal(oNamespace.Namespace.ID, out.Session.NamespaceId)
			ass.Len(out.Errors, 0)

			// check that code challenged & method are saved correctly
			{
				session, err := this.coreSession.Load(ctx, this.db, *out.Token)
				ass.NoError(err)
				ass.Equal("S256", session.CodeChallengeMethod)
				ass.Equal(
					"84d89877f0d4041efb6bf91a16f0248f2fd573e6af05c19f96bedb9f882f7882",
					session.CodeChallenge,
				)
			}

			// check that with outcome.Session we can generate JWT
			{
				signedString, err := this.SessionResolver.Jwt(ctx, out.Session)
				ass.NoError(err)
				ass.Contains(signedString, "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.")

				// check that JWT is valid
				{
					claims, err := this.SessionResolver.JwtValidation(signedString)
					ass.NoError(err)
					ass.NotNil(claims)
					ass.Equal(claims.SessionId(), out.Session.ID)
					ass.Equal(claims.UserId(), out.Session.UserId)
					ass.Equal(claims.NamespaceId(), out.Session.NamespaceId)
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
					NamespaceID: oNamespace.Namespace.ID,
					UserID:      oUser.User.ID,
				},
			})

			ass.NoError(err)
			ass.Equal(claim.KindOTLT, oGenerate.Session.Kind)

			// Use it
			{
				out, err := this.SessionCreate(ctx, &dto.SessionCreateInput{
					UseOTLT: &dto.SessionCreateUseOTLT{
						Token: *oGenerate.Token,
					},
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
	this := bean()

	// create user
	iUser := fUser.NewUserCreateInputFixture()

	// create namespace
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	oNamespace, err := this.namespace.NamespaceCreate(ctx, iNamespace)
	ass.NoError(err)

	// base input
	in := fixtures.SessionCreateInputFixtureUseCredentials(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	outcome, err := this.SessionCreate(ctx, in)
	ass.Error(err)
	ass.Nil(outcome)
	ass.Contains(err.Error(), "user not found")
}

func Test_Query(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := bean()

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := this.user.Resolvers.Mutation.UserCreate(ctx, iUser)

	ctx = context.WithValue(ctx, claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{Subject: oUser.User.ID},
		Kind:           claim.KindAuthenticated,
	})
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	oNamespace, _ := this.namespace.NamespaceCreate(ctx, iNamespace)
	in := fixtures.SessionCreateInputFixtureUseCredentials(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	outcome, err := this.SessionCreate(ctx, in)
	ass.NoError(err)

	// can load session without issue
	session, err := this.Session(ctx, *outcome.Token)
	ass.NoError(err)
	ass.Equal(session.NamespaceId, oNamespace.Namespace.ID)
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
		// …
	})
}

func Test_Archive(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := bean()

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := this.user.Resolvers.Mutation.UserCreate(ctx, iUser)
	ctx = context.WithValue(ctx, claim.ContextKey, &claim.Payload{
		StandardClaims: jwt.StandardClaims{Subject: oUser.User.ID},
		Kind:           claim.KindAuthenticated,
	})
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	oNamespace, _ := this.namespace.NamespaceCreate(ctx, iNamespace)
	in := fixtures.SessionCreateInputFixtureUseCredentials(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	sessionOutcome, err := this.SessionCreate(ctx, in)
	ass.NoError(err)

	// can archive session without issue
	{
		outcome, err := this.SessionArchive(ctx, *sessionOutcome.Token)
		ass.NoError(err)
		ass.Equal(outcome.Result, true)
	}

	// archive again -> should have error
	{
		outcome, err := this.SessionArchive(ctx, *sessionOutcome.Token)
		ass.NoError(err)
		ass.Equal(outcome.Result, false)
	}
}
