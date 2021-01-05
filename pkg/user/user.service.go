package user

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"bean/components/connect"
	"bean/components/util"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type UserService struct {
	bundle *Bundle
}

func (srv *UserService) Load(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	db := connect.ContextToDB(ctx)
	err := db.Where(&model.User{ID: id}).First(user).Error

	if nil != err {
		return nil, err
	}

	return user, nil
}

func (srv *UserService) Create(ctx context.Context, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
	db := connect.ContextToDB(ctx)

	if uint8(len(in.Emails.Secondary)) > srv.bundle.maxSecondaryEmailPerUser {
		return nil, errors.Wrap(
			util.ErrorInvalidArgument,
			fmt.Sprintf("too many secondary emails, limit is %d", srv.bundle.maxSecondaryEmailPerUser),
		)
	}

	// create base record
	obj := &model.User{
		ID:        srv.bundle.idr.ULID(),
		Version:   srv.bundle.idr.ULID(),
		AvatarURI: in.AvatarURI,
		IsActive:  in.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&obj).Error; nil != err {
		return nil, err
	}

	// create emails
	if err := srv.bundle.EmailService.CreateBulk(ctx, obj, in.Emails); nil != err {
		return nil, err
	}

	// save name object
	if err := srv.bundle.nameService.create(ctx, obj, in); nil != err {
		return nil, err
	}

	// save password
	if err := srv.bundle.PasswordService.create(ctx, obj, in.Password); nil != err {
		return nil, err
	}

	return &dto.UserMutationOutcome{User: obj}, nil
}

func (srv *UserService) Update(ctx context.Context, in dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
	tx := connect.ContextToDB(ctx)
	obj, err := srv.bundle.UserService.Load(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	// validate version
	if obj.Version != in.Version {
		errList := util.NewErrors(util.ErrorCodeConflict, []string{"version"}, "")

		return &dto.UserMutationOutcome{Errors: errList, User: nil}, nil
	}

	if err := srv.bundle.PasswordService.create(ctx, obj, in.Values.Password); nil != err {
		return nil, err
	}

	// bump new version
	obj.Version = srv.bundle.idr.ULID()
	if err := tx.Save(obj).Error; nil != err {
		return nil, err
	}

	return &dto.UserMutationOutcome{Errors: nil, User: obj}, nil
}
