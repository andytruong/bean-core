package access

import (
	"context"
	"path"
	"runtime"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/pkg/access/api/handler"
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
	config *Config,
) *AccessBean {
	bean := &AccessBean{
		config:    config.init(),
		logger:    logger,
		db:        db,
		id:        id,
		user:      bUser,
		namespace: bNamespace,
	}

	bean.SessionResolver = ModelResolver{
		logger: logger,
		bean:   bean,
		config: config,
	}

	return bean
}

type (
	AccessBean struct {
		config          *Config
		logger          *zap.Logger
		db              *gorm.DB
		id              *util.Identifier
		SessionResolver ModelResolver

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

func (this *AccessBean) SessionCreate(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	timeout, _ := time.ParseDuration("128h")
	if nil != this.config {
		timeout = this.config.SessionTimeout
	}

	hdl := handler.SessionCreateHandler{
		ID:             this.id,
		SessionTimeout: timeout,
		Namespace:      this.namespace,
	}

	txn := this.db.WithContext(ctx).Begin()
	outcome, err := hdl.Handle(ctx, txn, input)
	if nil != err {
		txn.Rollback()

		return nil, err
	}

	return outcome, txn.Commit().Error
}

func (this *AccessBean) SessionArchive(ctx context.Context, token string) (*dto.SessionDeleteOutcome, error) {
	session, err := this.Session(ctx, token)
	if nil != err {
		return &dto.SessionDeleteOutcome{
			Errors: util.NewErrors(util.ErrorCodeInput, []string{"token"}, err.Error()),
			Result: false,
		}, nil
	}

	hdl := handler.SessionDeleteHandler{
		ID: this.id,
		DB: this.db,
	}

	return hdl.Handle(ctx, session)
}

func (this AccessBean) Session(ctx context.Context, token string) (*model.Session, error) {
	hdl := handler.SessionLoadHandler{
		ID: this.id,
		DB: this.db,
	}

	return hdl.Handle(ctx, token)
}
