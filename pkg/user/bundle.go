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

func NewUserBundle(lgr *zap.Logger, idr *scalar.Identifier) *UserBundle {
	if err := util.NilPointerErrorValidate(lgr, idr); nil != err {
		panic(err)
	}

	this := &UserBundle{
		lgr:                      lgr,
		idr:                      idr,
		maxSecondaryEmailPerUser: 20,
	}

	this.Service = &UserService{bundle: this}
	this.nameService = &NameService{bundle: this}
	this.emailService = &EmailService{bundle: this}
	this.passwordService = &PasswordService{bundle: this}
	this.resolvers = newResolvers(this)

	return this
}

type UserBundle struct {
	module.AbstractBundle

	Service *UserService

	// Internal services
	lgr                      *zap.Logger
	idr                      *scalar.Identifier
	maxSecondaryEmailPerUser uint8
	resolvers                map[string]interface{}
	nameService              *NameService
	emailService             *EmailService
	passwordService          *PasswordService
}

func (UserBundle) Name() string {
	return "User"
}

func (bundle UserBundle) Migrate(ctx context.Context, driver string) error {
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

func (bundle *UserBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
