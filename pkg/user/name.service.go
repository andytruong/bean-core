package user

import (
	"context"

	"bean/components/connect"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type NameService struct {
	bundle *UserBundle
}

func (service *NameService) load(ctx context.Context, userId string) (*model.UserName, error) {
	name := &model.UserName{}
	err := connect.ContextToDB(ctx).Where(model.UserName{UserId: userId}).First(&name).Error
	if nil != err {
		return nil, err
	}

	return name, nil
}

func (service *NameService) create(ctx context.Context, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Name {
		name := model.UserName{
			ID:            service.bundle.idr.MustULID(),
			UserId:        user.ID,
			FirstName:     input.Name.FirstName,
			LastName:      input.Name.LastName,
			PreferredName: input.Name.PreferredName,
		}

		if err := connect.ContextToDB(ctx).Create(name).Error; nil != err {
			return err
		}
	}

	return nil
}
