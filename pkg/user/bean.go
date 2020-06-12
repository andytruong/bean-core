package user

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"bean/pkg/user/api/handler"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewUserBean(db *gorm.DB, logger *zap.Logger, id *util.Identifier) *UserBean {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		panic(err)
	}

	this := &UserBean{
		logger:                   logger,
		db:                       db,
		id:                       id,
		maxSecondaryEmailPerUser: 20,
	}

	this.Email = &UserBeanEmail{bean: this}

	return this
}

type UserBean struct {
	logger                   *zap.Logger
	db                       *gorm.DB
	id                       *util.Identifier
	maxSecondaryEmailPerUser uint8
	Email                    *UserBeanEmail
}

func (this UserBean) Dependencies() []util.Bean {
	return nil
}

func (this UserBean) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "user",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this UserBean) Verified(ctx context.Context, obj *model.UserEmail) (bool, error) {
	return obj.IsVerified, nil
}

func (this *UserBean) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
	hdl := handler.UserCreateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})

	if uint8(len(input.Emails.Secondary)) > this.maxSecondaryEmailPerUser {
		return nil, errors.Wrap(
			util.ErrorInvalidArgument,
			fmt.Sprintf("too many secondary emails, limit is %d", this.maxSecondaryEmailPerUser),
		)
	}

	if outcome, err := hdl.Handle(txn, input); nil != err {
		txn.Rollback()

		return nil, err
	} else {
		txn.Commit()

		return outcome, nil
	}
}

func (this *UserBean) User(ctx context.Context, id string) (*model.User, error) {
	hdl := handler.UserLoadHandler{}

	return hdl.Load(this.db, id)
}

func (this UserBean) Name(ctx context.Context, user *model.User) (*model.UserName, error) {
	name := model.UserName{}
	err := this.db.
		Where(model.UserName{UserId: user.ID}).
		First(&name).
		Error

	if nil != err {
		return nil, errors.Wrap(util.ErrorQuery, err.Error())
	}

	return &name, nil
}

func (this UserBean) Emails(ctx context.Context, obj *model.User) (*model.UserEmails, error) {
	return this.Email.List(ctx, obj)
}

func (this UserBean) UserUpdate(ctx context.Context, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
	user, err := this.User(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	hdl := handler.UserUpdateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})

	if outcome, err := hdl.Handle(txn, user, input); nil != err {
		txn.Rollback()

		return nil, err
	} else {
		txn.Commit()

		return outcome, nil
	}
}
