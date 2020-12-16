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
	db *gorm.DB,
	id *unique.Identifier,
	logger *zap.Logger,
	userBundle *user.UserBundle,
	spaceBundle *space.SpaceBundle,
	config *AccessConfiguration,
) *AccessBundle {
	this := &AccessBundle{
		config:      config.init(),
		logger:      logger,
		db:          db,
		id:          id,
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

		config         *AccessConfiguration
		logger         *zap.Logger
		db             *gorm.DB
		id             *unique.Identifier
		sessionService *SessionService
		JwtService     *JwtService

		// depends on userBundle bundle
		userBundle  *user.UserBundle
		spaceBundle *space.SpaceBundle
		resolvers   map[string]interface{}
	}
)

func (this AccessBundle) Dependencies() []module.Bundle {
	return []module.Bundle{
		this.userBundle,
		this.spaceBundle,
	}
}

func (this AccessBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "access",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this AccessBundle) GraphqlResolver() map[string]interface{} {
	return this.resolvers
}
