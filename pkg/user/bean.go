package user

import (
	"context"
	"path"
	"runtime"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
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

	this.Core = &Core{bean: this}
	this.CoreName = &CoreName{bean: this}
	this.CoreEmail = &CoreEmail{bean: this}
	this.CorePassword = &CorePassword{bean: this}

	return this
}

type UserBean struct {
	logger                   *zap.Logger
	db                       *gorm.DB
	id                       *util.Identifier
	maxSecondaryEmailPerUser uint8
	Core                     *Core
	CoreName                 *CoreName
	CoreEmail                *CoreEmail
	CorePassword             *CorePassword
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

func (this *UserBean) UserCreate(ctx context.Context, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
	var err error
	var out *dto.UserMutationOutcome

	err = connect.Transaction(ctx, this.db, func(tx *gorm.DB) error {
		out, err = this.Core.Create(tx, in)

		return err
	})

	return out, err
}

func (this *UserBean) User(ctx context.Context, id string) (*model.User, error) {
	return this.Core.Load(this.db, id)
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
	return this.CoreEmail.List(ctx, obj)
}

func (this UserBean) UserUpdate(ctx context.Context, input dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
	var err error
	var out *dto.UserMutationOutcome

	err = connect.Transaction(ctx, this.db, func(tx *gorm.DB) error {
		out, err = this.Core.Update(tx, input)

		return err
	})

	return out, err
}
