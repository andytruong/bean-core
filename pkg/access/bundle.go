package access

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/space"
	"bean/pkg/user"
)

func NewAccessBundle(
	idr *unique.Identifier,
	logger *zap.Logger,
	userBundle *user.UserBundle,
	spaceBundle *space.SpaceBundle,
	cnf *AccessConfiguration,
) *AccessBundle {
	this := &AccessBundle{
		cnf:         cnf.init(),
		logger:      logger,
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
		logger         *zap.Logger
		idr            *unique.Identifier
		sessionService *SessionService
		JwtService     *JwtService

		// depends on userBundle bundle
		userBundle  *user.UserBundle
		spaceBundle *space.SpaceBundle
		resolvers   map[string]interface{}
	}
)

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
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "access",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle AccessBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
