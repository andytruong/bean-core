package user

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bean/components/util"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type UserService struct {
	bundle *UserBundle
}

func (service *UserService) Load(db *gorm.DB, id string) (*model.User, error) {
	user := &model.User{}
	err := db.Where(&model.User{ID: id}).First(user).Error

	if nil != err {
		return nil, err
	}

	return user, nil
}

func (service *UserService) Create(tx *gorm.DB, in *dto.UserCreateInput) (*dto.UserMutationOutcome, error) {
	if uint8(len(in.Emails.Secondary)) > service.bundle.maxSecondaryEmailPerUser {
		return nil, errors.Wrap(
			util.ErrorInvalidArgument,
			fmt.Sprintf("too many secondary emails, limit is %d", service.bundle.maxSecondaryEmailPerUser),
		)
	}

	// create base record
	obj := &model.User{
		ID:        service.bundle.id.MustULID(),
		Version:   service.bundle.id.MustULID(),
		AvatarURI: in.AvatarURI,
		IsActive:  in.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tx.Create(&obj).Error; nil != err {
		return nil, err
	}

	// create emails
	if err := service.bundle.emailService.CreateBulk(tx, obj, in.Emails); nil != err {
		return nil, err
	}

	// save name object
	if err := service.bundle.nameService.create(tx, obj, in); nil != err {
		return nil, err
	}

	// save password
	if err := service.bundle.passwordService.create(tx, obj, in.Password); nil != err {
		return nil, err
	}

	return &dto.UserMutationOutcome{User: obj}, nil
}

func (service *UserService) Update(tx *gorm.DB, in dto.UserUpdateInput) (*dto.UserMutationOutcome, error) {
	obj, err := service.bundle.Service.Load(tx, in.ID)
	if err != nil {
		return nil, err
	}

	// validate version
	if obj.Version != in.Version {
		errList := util.NewErrors(util.ErrorCodeConflict, []string{"version"}, "")

		return &dto.UserMutationOutcome{Errors: errList, User: nil}, nil
	}

	if err := service.bundle.passwordService.create(tx, obj, in.Values.Password); nil != err {
		return nil, err
	}

	// bump new version
	obj.Version = service.bundle.id.MustULID()
	if err := tx.Save(obj).Error; nil != err {
		return nil, err
	}

	return &dto.UserMutationOutcome{Errors: nil, User: obj}, nil
}
