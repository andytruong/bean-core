package user

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"

	"bean/pkg/user/api"
	"bean/pkg/user/dto"
	"bean/pkg/util"
)

type (
	UserMutationResolver struct {
		db *gorm.DB
		id *util.Identifier
	}
)

// TODO: Work in progress
// TODO: validate avatar URI
func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	ctl := api.UserCreateAPI{ID: this.id}
	tx := this.db.BeginTx(ctx, &sql.TxOptions{})

	if outcome, err := ctl.Create(tx, input); nil != err {
		tx.Rollback()

		return nil, err
	} else {
		tx.Commit()

		return outcome, nil
	}
}
