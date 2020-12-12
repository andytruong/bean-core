package access

import (
	"context"
	"path"
	"runtime"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	
	"bean/components/claim"
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
	
	return this
}

type (
	AccessBundle struct {
		config          *AccessConfiguration
		logger          *zap.Logger
		db              *gorm.DB
		id              *unique.Identifier
		SessionResolver ModelResolver
		sessionService  *SessionService
		
		// depends on userBundle bundle
		userBundle  *user.UserBundle
		spaceBundle *space.SpaceBundle
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

func (this *AccessBundle) SessionCreate(ctx context.Context, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	txn := this.db.WithContext(ctx).Begin()
	outcome, err := this.sessionService.Create(txn, in)
	if nil != err {
		txn.Rollback()
		
		return nil, err
	}
	
	return outcome, txn.Commit().Error
}

func (this *AccessBundle) SessionArchive(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
	claims := claim.ContextToPayload(ctx)
	if nil == claims {
		return nil, util.ErrorAuthRequired
	}
	
	tx := this.db.WithContext(ctx).Begin()
	
	// load session
	sess, err := this.sessionService.load(ctx, tx, claims.SessionId())
	if nil != err {
		return &dto.SessionArchiveOutcome{
			Errors: util.NewErrors(util.ErrorCodeInput, []string{"token"}, err.Error()),
			Result: false,
		}, nil
	}
	
	// delete it
	{
		out, err := this.sessionService.Delete(tx, sess)
		if nil != err {
			tx.Rollback()
			
			return nil, err
		}
		
		return out, tx.Commit().Error
	}
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
