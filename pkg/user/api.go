package user

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/dto"
	"bean/pkg/user/model"
	"bean/pkg/user/service"
	"bean/pkg/util"
)

func NewUserMutationResolver(db *gorm.DB, id *util.Identifier) (*UserMutationResolver, error) {
	if err := util.NilPointerErrorValidate(db, id); nil != err {
		return nil, err
	}

	return &UserMutationResolver{
		db: db,
		id: id,
	}, nil
}

type (
	UserMutationResolver struct {
		db *gorm.DB
		id *util.Identifier
	}

	UserQueryResolver struct {
		db *gorm.DB
	}
)

func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	sv := service.UserCreateAPI{ID: this.id}
	tx := this.db.BeginTx(ctx, &sql.TxOptions{})

	if outcome, err := sv.Create(tx, input); nil != err {
		tx.Rollback()

		return nil, err
	} else {
		tx.Commit()

		return outcome, nil
	}
}

func (this *UserQueryResolver) User(ctx context.Context, id string) (*model.User, error) {
	sv := service.UserQueryAPI{}

	return sv.Load(this.db, id)
}
