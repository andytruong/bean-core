package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

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
		db:                       db,
		id:                       id,
		maxSecondaryEmailPerUser: 20,
	}, nil
}

type (
	UserMutationResolver struct {
		db                       *gorm.DB
		id                       *util.Identifier
		maxSecondaryEmailPerUser uint8
	}

	UserQueryResolver struct {
		db *gorm.DB
	}
)

func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	sv := service.UserCreateAPI{ID: this.id}
	tx := this.db.BeginTx(ctx, &sql.TxOptions{})

	if uint8(len(input.Emails.Secondary)) > this.maxSecondaryEmailPerUser {
		return nil, errors.Wrap(
			util.ErrorInvalidArgument,
			fmt.Sprintf("too many secondary emails, limit is %d", this.maxSecondaryEmailPerUser),
		)
	}

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
