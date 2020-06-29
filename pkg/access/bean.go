package access

import (
	"context"
	"path"
	"runtime"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/namespace"
	"bean/pkg/user"
	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewAccessBean(
	db *gorm.DB,
	id *util.Identifier,
	logger *zap.Logger,
	bUser *user.UserBean,
	bNamespace *namespace.NamespaceBean,
	genetic *Genetic,
) *AccessBean {
	this := &AccessBean{
		genetic:   genetic.init(),
		logger:    logger,
		db:        db,
		id:        id,
		user:      bUser,
		namespace: bNamespace,
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
		id              *util.Identifier
		SessionResolver ModelResolver
		coreSession     *CoreSession

		// depends on user bean
		user      *user.UserBean
		namespace *namespace.NamespaceBean
	}
)

func (this AccessBean) Dependencies() []util.Bean {
	return []util.Bean{
		this.user,
		this.namespace,
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

func (this *AccessBean) SessionArchive(ctx context.Context, token string) (*dto.SessionDeleteOutcome, error) {
	tx := this.db.WithContext(ctx).Begin()

	// load session
	sess, err := this.coreSession.Load(ctx, tx, token)
	if nil != err {
		return &dto.SessionDeleteOutcome{
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
	return this.coreSession.Load(ctx, this.db, token)
}
