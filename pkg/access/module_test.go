package access

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"bean/pkg/access/api/fixtures"
	"bean/pkg/namespace"
	fNamespace "bean/pkg/namespace/api/fixtures"
	"bean/pkg/user"
	fUser "bean/pkg/user/api/fixtures"
	"bean/pkg/util"
)

func module() *AccessModule {
	configRaw, err := util.ParseFile("../../config.yaml")
	if nil != err {
		panic(err)
	}

	config := &struct {
		Modules struct {
			Access *Config `yaml:"access"`
		} `yaml:"modules"`
	}{}

	{
		err := yaml.Unmarshal(configRaw, &config)
		if nil != err {
			panic(err)
		}
	}

	db := util.MockDatabase()
	logger := util.MockLogger()
	id := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, id)
	mNamespace := namespace.NewNamespaceModule(db, logger, id, mUser)
	module := NewAccessModule(db, id, logger, mUser, mNamespace, config.Modules.Access)
	util.MockInstall(mUser, db)
	util.MockInstall(mNamespace, db)
	util.MockInstall(module, db)

	return module
}

func Test_Config(t *testing.T) {
	ass := assert.New(t)
	this := module()
	key, err := this.config.GetSignKey()
	ass.NoError(err)
	ass.NotNil(key)
}

func Test_Create(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := module()

	// create user
	iUser := fUser.NewUserCreateInputFixture()
	oUser, err := this.userModule.UserCreate(ctx, iUser)
	ass.NoError(err)

	// create namespace
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	iNamespace.Context.UserID = oUser.User.ID
	oNamespace, err := this.namespaceModule.NamespaceCreate(ctx, iNamespace)
	ass.NoError(err)

	t.Run("test email inactive", func(t *testing.T) {
		input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
		input.Email = iUser.Emails.Secondary[1].Value
		_, err := this.SessionCreate(ctx, input)
		ass.Equal(err.Error(), "user not found")
	})

	t.Run("test password unmatched", func(t *testing.T) {
		input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
		input.HashedPassword = "invalid-password"
		outcome, err := this.SessionCreate(ctx, input)

		ass.NoError(err)
		ass.Equal(util.ErrorCodeInput, *outcome.Errors[0].Code)
		ass.Equal(outcome.Errors[0].Message, "invalid password")
		ass.Equal(outcome.Errors[0].Fields, []string{"input.namespaceId"})
	})

	t.Run("test ok", func(t *testing.T) {
		input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)
		outcome, err := this.SessionCreate(ctx, input)
		ass.NoError(err)
		ass.Equal(oUser.User.ID, outcome.Session.UserId)
		ass.Equal(oNamespace.Namespace.ID, outcome.Session.NamespaceId)
		ass.Len(outcome.Errors, 0)

		{
			// check that with outcome.Session we can generate JWT
			signedString, err := this.SessionResolver.Jwt(ctx, outcome.Session)
			ass.NoError(err)
			ass.Contains(signedString, "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.")

			{
				// check that JWT is valid
				claims, err := this.SessionResolver.JwtValidation(signedString)
				ass.NoError(err)
				ass.NotNil(claims)
				ass.Equal(claims.SessionId(), outcome.Session.ID)
				ass.Equal(claims.UserId(), outcome.Session.UserId)
				ass.Equal(claims.NamespaceId(), outcome.Session.NamespaceId)
			}
		}
	})
}

func Test_SessionCreate_MembershipNotFound(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := module()

	// create user
	iUser := fUser.NewUserCreateInputFixture()

	// create namespace
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	oNamespace, err := this.namespaceModule.NamespaceCreate(ctx, iNamespace)
	ass.NoError(err)

	// base input
	input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	outcome, err := this.SessionCreate(ctx, input)
	ass.Error(err)
	ass.Nil(outcome)
	ass.Contains(err.Error(), "user not found")
}

func Test_Query(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := module()

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := this.userModule.UserCreate(ctx, iUser)
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	iNamespace.Context.UserID = oUser.User.ID
	oNamespace, _ := this.namespaceModule.NamespaceCreate(ctx, iNamespace)
	input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	outcome, err := this.SessionCreate(ctx, input)
	ass.NoError(err)

	// can load session without issue
	session, err := this.Session(ctx, *outcome.Token)
	ass.NoError(err)
	ass.Equal(session.NamespaceId, oNamespace.Namespace.ID)
	ass.Equal(session.UserId, oUser.User.ID)

	{
		// change session expiration time
		oneMinDuration, _ := time.ParseDuration("129h")
		session.ExpiredAt = session.ExpiredAt.Add(-1 * oneMinDuration)
		err := this.db.Table("access_session").Save(&session).Error
		ass.NoError(err)

		// load again -> error: session expired
		_, err = this.Session(ctx, *outcome.Token)
		ass.Error(err)
		ass.Equal(err.Error(), "session expired")
	}
}

func Test_Archive(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := module()

	iUser := fUser.NewUserCreateInputFixture()
	oUser, _ := this.userModule.UserCreate(ctx, iUser)
	iNamespace := fNamespace.NamespaceCreateInputFixture(false)
	iNamespace.Context.UserID = oUser.User.ID
	oNamespace, _ := this.namespaceModule.NamespaceCreate(ctx, iNamespace)
	input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, string(iUser.Emails.Secondary[0].Value), iUser.Password.HashedValue)

	sessionOutcome, err := this.SessionCreate(ctx, input)
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
