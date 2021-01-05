package user

import (
	"context"
	"path"
	"runtime"

	"go.uber.org/zap"

	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/components/util"
)

const (
	ErrorUserNotFound    = util.Err("user not found")
	ErrorEmailInactive   = util.Err("email address is not active")
	ErrorInvalidPassword = util.Err("invalid password")
)

func NewUserBundle(lgr *zap.Logger, idr *scalar.Identifier) *Bundle {
	if err := util.NilPointerErrorValidate(lgr, idr); nil != err {
		panic(err)
	}

	this := &Bundle{
		lgr:                      lgr,
		idr:                      idr,
		maxSecondaryEmailPerUser: 20,
	}

	this.UserService = &UserService{bundle: this}
	this.EmailService = &EmailService{bundle: this}
	this.PasswordService = &PasswordService{bundle: this}
	this.nameService = &NameService{bundle: this}
	this.resolvers = newResolvers(this)

	return this
}

type Bundle struct {
	module.AbstractBundle

	UserService     *UserService
	EmailService    *EmailService
	PasswordService *PasswordService

	// Internal services
	lgr                      *zap.Logger
	idr                      *scalar.Identifier
	maxSecondaryEmailPerUser uint8
	resolvers                map[string]interface{}
	nameService              *NameService
}

func (Bundle) Name() string {
	return "User"
}

func (bundle Bundle) Migrate(ctx context.Context, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := connect.Runner{
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "user",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run(ctx)
}

func (bundle *Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
