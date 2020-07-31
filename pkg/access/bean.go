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

func NewAccessBean(
	db *gorm.DB,
	id *unique.Identifier,
	logger *zap.Logger,
	bUser *user.UserBean,
	bSpace *space.SpaceBean,
	genetic *Genetic,
) *AccessBean {
	this := &AccessBean{
		genetic: genetic.init(),
		logger:  logger,
		db:      db,
		id:      id,
		user:    bUser,
		space:   bSpace,
	}

	this.SessionResolver = ModelResolver{
		logger: logger,
		bean:   this,
		config: genetic,
	}

	this.coreSession = &CoreSession{bean: this}

	return this
}

type (
	AccessBean struct {
		genetic         *Genetic
		logger          *zap.Logger
		db              *gorm.DB
		id              *unique.Identifier
		SessionResolver ModelResolver
		coreSession     *CoreSession

		// depends on user bean
		user  *user.UserBean
		space *space.SpaceBean
	}
)

func (this AccessBean) Dependencies() []module.Bean {
	return []module.Bean{
		this.user,
		this.space,
	}
}

func (this AccessBean) Migrate(tx *gorm.DB, driver string) error {
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

func (this *AccessBean) SessionCreate(ctx context.Context, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	txn := this.db.WithContext(ctx).Begin()
	outcome, err := this.coreSession.Create(txn, in)
	if nil != err {
		txn.Rollback()

		return nil, err
	}

	return outcome, txn.Commit().Error
}

func (this *AccessBean) SessionArchive(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
	claims := claim.ContextToPayload(ctx)
	if nil == claims {
		return nil, util.ErrorAuthRequired
	}

	tx := this.db.WithContext(ctx).Begin()

	// load session
	sess, err := this.coreSession.load(ctx, tx, claims.SessionId())
	if nil != err {
		return &dto.SessionArchiveOutcome{
			Errors: util.NewErrors(util.ErrorCodeInput, []string{"token"}, err.Error()),
			Result: false,
		}, nil
	}

	// delete it
	{
		out, err := this.coreSession.Delete(tx, sess)
		if nil != err {
			tx.Rollback()

			return nil, err
		}

		return out, tx.Commit().Error
	}
}

func (this AccessBean) Session(ctx context.Context, token string) (*model.Session, error) {
	return this.coreSession.LoadByToken(ctx, this.db, token)
}

func (this AccessBean) Sign(claims jwt.Claims) (string, error) {
	key, err := this.genetic.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}

	return jwt.
		NewWithClaims(this.genetic.signMethod(), claims).
		SignedString(key)
}
