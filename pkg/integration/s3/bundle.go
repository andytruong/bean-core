package s3

import (
	"path"
	"runtime"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
)

func NewS3Integration(
	db *gorm.DB,
	id *unique.Identifier,
	logger *zap.Logger,
	conf *S3Configuration,
) *S3IntegrationBundle {
	this := &S3IntegrationBundle{
		db:     db,
		id:     id,
		logger: logger,
		config: conf,
	}
	
	this.AppService = &ApplicationService{
		bundle:   this,
		Resolver: &ApplicationResolver{bundle: this},
	}
	this.credentialService = &credentialService{bundle: this}
	this.policyService = &policyService{bundle: this}
	
	return this
}

type S3IntegrationBundle struct {
	module.AbstractBundle
	
	db     *gorm.DB
	id     *unique.Identifier
	logger *zap.Logger
	config *S3Configuration
	
	AppService        *ApplicationService
	credentialService *credentialService
	policyService     *policyService
}

func (this S3IntegrationBundle) Migrate(tx *gorm.DB, driver string) error {
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

func (this S3IntegrationBundle) GraphqlResolver() map[string]interface{} {
	// TODO: Singleton
	
	return newGraphqlResolver()
}
