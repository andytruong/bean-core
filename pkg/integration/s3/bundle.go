package s3

import (
	"path"
	"runtime"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/scalar"
	"bean/pkg/app"
)

func NewS3Integration(
	idr *scalar.Identifier,
	logger *zap.Logger,
	conf *S3Configuration,
	appBundle *app.AppBundle,
) *S3Bundle {
	this := &S3Bundle{
		idr:    idr,
		logger: logger,
		cnf:    conf,
	}

	this.appBundle = appBundle
	this.AppService = &ApplicationService{bundle: this}
	this.credentialService = &credentialService{bundle: this}
	this.policyService = &policyService{bundle: this}
	this.resolvers = newResolvers(this)

	return this
}

type S3Bundle struct {
	module.AbstractBundle

	idr               *scalar.Identifier
	logger            *zap.Logger
	cnf               *S3Configuration
	appBundle         *app.AppBundle
	AppService        *ApplicationService
	credentialService *credentialService
	policyService     *policyService
	resolvers         map[string]interface{}
}

func (bundle S3Bundle) Dependencies() []module.Bundle {
	return []module.Bundle{bundle.appBundle}
}

func (bundle S3Bundle) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: bundle.logger,
		Driver: driver,
		Bundle: "integration.s3",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (bundle *S3Bundle) GraphqlResolver() map[string]interface{} {
	return bundle.resolvers
}
