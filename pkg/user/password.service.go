package user

import (
	"context"
	"time"

	"bean/components/connect"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type PasswordService struct {
	bundle *UserBundle
}

func (service *PasswordService) create(ctx context.Context, user *model.User, in *dto.UserPasswordInput) error {
	if nil == in {
		return nil
	}

	db := connect.ContextToDB(ctx)
	pass := &model.UserPassword{
		ID:          service.bundle.idr.MustULID(),
		UserId:      user.ID,
		HashedValue: in.HashedValue,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
	}

	{
		err := db.Create(pass).Error
		if nil != err {
			return err
		}
	}

	// set other passwords to inactive
	return db.
		Where("user_id == ?", pass.UserId).
		Where("id != ?", pass.ID).
		Updates(model.UserPassword{IsActive: false}).
		Error
}
