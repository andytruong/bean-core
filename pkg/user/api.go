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

	UserModelResolver struct {
		db *gorm.DB
	}

	UserEmailResolver struct {
		db *gorm.DB
	}
)

func (this UserEmailResolver) Verified(ctx context.Context, obj *model.UserEmail) (bool, error) {
	return obj.IsVerified, nil
}

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

// TODO: dataloader
func (this UserModelResolver) Name(ctx context.Context, user *model.User) (*model.UserName, error) {
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

// TODO: dataloader
// TODO: need a better resolver, we not always load secondary emails.
func (this UserModelResolver) Emails(ctx context.Context, user *model.User) (*model.UserEmails, error) {
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
