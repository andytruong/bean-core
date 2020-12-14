package access

import (
	"context"
	"path"
	"runtime"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/module"
	"bean/components/module/migrate"
	"bean/components/unique"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/space"
	"bean/pkg/user"
	"bean/pkg/util"
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
	
	this.SessionResolver = ModelResolver{
		logger: logger,
		bundle: this,
		config: config,
	}
	
	this.sessionService = &SessionService{bundle: this}
	this.resolvers = this.newResolves()
	
	return this
}

type (
	AccessBundle struct {
		module.AbstractBundle
		
		config          *AccessConfiguration
		logger          *zap.Logger
		db              *gorm.DB
		id              *unique.Identifier
		SessionResolver ModelResolver
		sessionService  *SessionService
		
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

func (this AccessBundle) Session(ctx context.Context, token string) (*model.Session, error) {
	return this.sessionService.LoadByToken(ctx, this.db, token)
}

func (this AccessBundle) Sign(claims jwt.Claims) (string, error) {
	key, err := this.config.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}
	
	return jwt.
		NewWithClaims(this.config.signMethod(), claims).
		SignedString(key)
}

func (this AccessBundle) GraphqlResolver() map[string]interface{} {
	return this.resolvers
}

func (this AccessBundle) newResolves() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{},
		"Mutation": map[string]interface{}{
			"SessionCreate": func(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
				return this.SessionCreate(ctx, input)
			},
			"SessionArchive": func(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
				return this.SessionArchive(ctx)
			},
		},
	}
}
