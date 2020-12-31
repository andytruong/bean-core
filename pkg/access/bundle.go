package access

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
	"bean/pkg/space"
	"bean/pkg/user"
)

func NewAccessBundle(
	idr *scalar.Identifier,
	lgr *zap.Logger,
	userBundle *user.UserBundle,
	spaceBundle *space.SpaceBundle,
	cnf *AccessConfiguration,
) *AccessBundle {
	this := &AccessBundle{
		cnf:         cnf.init(),
		lgr:         lgr,
		idr:         idr,
		userBundle:  userBundle,
		spaceBundle: spaceBundle,
	}

	this.sessionService = &SessionService{bundle: this}
	this.JwtService = &JwtService{bundle: this}
	this.resolvers = this.newResolves()

	return this
}

type (
	AccessBundle struct {
		module.AbstractBundle

		cnf            *AccessConfiguration
		lgr            *zap.Logger
		idr            *scalar.Identifier
		sessionService *SessionService
		JwtService     *JwtService

		// depends on userBundle bundle
		userBundle  *user.UserBundle
		spaceBundle *space.SpaceBundle
		resolvers   map[string]interface{}
	}
)

func (AccessBundle) Name() string {
	return "Access"
}

func (bundle AccessBundle) Dependencies() []module.Bundle {
	return []module.Bundle{
		bundle.userBundle,
		bundle.spaceBundle,
	}
}

func (bundle AccessBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.lgr,
		Driver: driver,
		Bundle: "access",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle AccessBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
