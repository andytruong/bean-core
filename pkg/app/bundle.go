package app

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/space"
)

func NewApplicationBundle(
	con *gorm.DB,
	idr *unique.Identifier,
	logger *zap.Logger,
) (*AppBundle, error) {
	bundle := &AppBundle{
		con:         con,
		idr:         idr,
		logger:      logger,
		spaceBundle: nil,
	}

	bundle.resolvers = bundle.newResolvers()

	return bundle, nil
}

type AppBundle struct {
	module.AbstractBundle

	spaceBundle *space.SpaceBundle
	con         *gorm.DB
	idr         *unique.Identifier
	logger      *zap.Logger
	resolvers   map[string]interface{}
}

func (bundle AppBundle) Dependencies() []module.Bundle {
	if nil != bundle.spaceBundle {
		return []module.Bundle{bundle.spaceBundle}
	}

	return nil
}

func (bundle AppBundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "bundle",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle *AppBundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
