package s3

import (
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewS3Integration(
	db *gorm.DB,
	id *util.Identifier,
	logger *zap.Logger,
	genetic *Genetic,
) *S3IntegrationBean {
	this := &S3IntegrationBean{
		db:      db,
		id:      id,
		logger:  logger,
		genetic: genetic,
	}

	this.CoreApp = &CoreApplication{
		bean:     this,
		Resolver: &ApplicationResolver{bean: this},
	}
	this.coreCredentials = &coreCredentials{bean: this}
	this.corePolicy = &corePolicy{bean: this}

	return this
}

type S3IntegrationBean struct {
	db      *gorm.DB
	id      *util.Identifier
	logger  *zap.Logger
	genetic *Genetic

	CoreApp         *CoreApplication
	coreCredentials *coreCredentials
	corePolicy      *corePolicy
}

func (this S3IntegrationBean) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "integration.s3",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this S3IntegrationBean) Dependencies() []util.Bean {
	return nil
}
