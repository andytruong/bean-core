package access

import (
	"context"
	"database/sql"
	"path"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/access/api/handler"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/namespace"
	"bean/pkg/user"
	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewAccessModule(
	db *gorm.DB,
	id *util.Identifier,
	logger *zap.Logger,
	userModule *user.UserModule,
	namespaceModule *namespace.NamespaceModule,
	config *Config,
) *AccessModule {
	module := &AccessModule{
		config:          config,
		logger:          logger,
		db:              db,
		id:              id,
		userModule:      userModule,
		namespaceModule: namespaceModule,
	}

	module.SessionResolver = ModelResolver{module: module}

	return module
}

type (
	AccessModule struct {
		config          *Config
		logger          *zap.Logger
		db              *gorm.DB
		id              *util.Identifier
		SessionResolver ModelResolver

		// depends on user module
		userModule      *user.UserModule
		namespaceModule *namespace.NamespaceModule
	}

	Config struct {
		SessionTimeout time.Duration `yaml:"sessionTimeout"`
	}
)

func (this AccessModule) Dependencies() []util.Module {
	return []util.Module{
		this.userModule,
		this.namespaceModule,
	}
}

func (this AccessModule) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Module: "access",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this *AccessModule) SessionCreate(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	timeout, _ := time.ParseDuration("128h")
	if nil != this.config {
		timeout = this.config.SessionTimeout
	}

	hdl := handler.SessionCreateHandler{
		ID:             this.id,
		SessionTimeout: timeout,
	}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})
	outcome, err := hdl.Handle(ctx, txn, input)
	if nil != err {
		txn.Rollback()

		return nil, err
	}

	return outcome, txn.Commit().Error
}

func (this *AccessModule) SessionArchive(ctx context.Context, token string) (*dto.SessionDeleteOutcome, error) {
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

func (this AccessModule) Session(ctx context.Context, token string) (*model.Session, error) {
	hdl := handler.SessionLoadHandler{
		ID: this.id,
		DB: this.db,
	}

	return hdl.Handle(ctx, token)
}

func (this AccessModule) Jwt(ctx context.Context, obj *dto.SessionCreateOutcome) (*string, error) {
	claims := util.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "access",
			Id:        obj.Session.ID,
			IssuedAt:  time.Now().UnixNano(),
			ExpiresAt: time.Now().Add(5 * time.Minute).UnixNano(),
			Subject:   obj.Session.UserId,
			Audience:  obj.Session.NamespaceId,
		},
	}

	mySigningKey := []byte("AllYourBase")
	signedString, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(mySigningKey)

	return util.NilString(signedString), err
}
