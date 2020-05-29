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

func NewUserModule(db *gorm.DB, logger *zap.Logger, id *util.Identifier) *UserModule {
	if err := util.NilPointerErrorValidate(db, logger, id); nil != err {
		panic(err)
	}

	return &UserModule{
		logger:                   logger,
		db:                       db,
		id:                       id,
		maxSecondaryEmailPerUser: 20,
	}
}

type UserModule struct {
	logger                   *zap.Logger
	db                       *gorm.DB
	id                       *util.Identifier
	maxSecondaryEmailPerUser uint8
}

func (this UserModule) Dependencies() []util.Module {
	return nil
}

func (this UserModule) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Module: "user",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this UserModule) Verified(ctx context.Context, obj *model.UserEmail) (bool, error) {
	return obj.IsVerified, nil
}

func (this *UserModule) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	hdl := handler.UserCreateHandler{ID: this.id}
	txn := this.db.BeginTx(ctx, &sql.TxOptions{})

	if uint8(len(input.Emails.Secondary)) > this.maxSecondaryEmailPerUser {
		return nil, errors.Wrap(
			util.ErrorInvalidArgument,
			fmt.Sprintf("too many secondary emails, limit is %d", this.maxSecondaryEmailPerUser),
		)
	}

	if outcome, err := hdl.Create(txn, input); nil != err {
		txn.Rollback()

		return nil, err
	} else {
		txn.Commit()

		return outcome, nil
	}
}

func (this *UserModule) User(ctx context.Context, id string) (*model.User, error) {
	hdl := handler.UserLoadHandler{}

	return hdl.Load(this.db, id)
}

func (this UserModule) Name(ctx context.Context, user *model.User) (*model.UserName, error) {
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

// TODO: need a better resolver, we not always load secondary emails.
func (this UserModule) Emails(ctx context.Context, user *model.User) (*model.UserEmails, error) {
	emails := &model.UserEmails{}

	var rows []*model.UserEmail
	err := this.db.
		Raw(`
			     SELECT *, 1 AS is_verified FROM user_emails            WHERE user_id = ?
		   UNION SELECT *, 0 AS is_verified FROM user_unverified_emails WHERE user_id = ?
		`, user.ID, user.ID).
		Find(&rows).
		Error

	if nil != err {
		return nil, err
	} else {
		for _, row := range rows {
			if row.IsPrimary {
				emails.Primary = row
			} else {
				emails.Secondary = append(emails.Secondary, row)
			}
		}
	}

	return emails, nil
}
