package user

import (
	"context"
	"time"

	"bean/components/connect"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type PasswordService struct {
	bundle *Bundle
}

func (srv PasswordService) ValidPassword(ctx context.Context, userId string, hashedValue string) (bool, error) {
	count := int64(0)
	err := connect.
		ContextToDB(ctx).
		Model(&model.UserPassword{}).
		Where("user_id = ? AND hashed_value = ? AND is_active = ?", userId, hashedValue, true).
		Count(&count).
		Error

	if nil != err {
		return false, err
	}

	return count == 1, nil
}

func (srv *PasswordService) create(ctx context.Context, user *model.User, in *dto.UserPasswordInput) error {
	if nil == in {
		return nil
	}

	db := connect.ContextToDB(ctx)
	pass := &model.UserPassword{
		ID:          srv.bundle.idr.ULID(),
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
