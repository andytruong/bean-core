package user

import (
	"context"

	"bean/components/connect"
	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type NameService struct {
	bundle *Bundle
}

func (srv *NameService) load(ctx context.Context, userId string) (*model.UserName, error) {
	name := &model.UserName{}
	err := connect.DB(ctx).Where(model.UserName{UserId: userId}).Take(&name).Error
	if nil != err {
		return nil, err
	}

	return name, nil
}

func (srv *NameService) create(ctx context.Context, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Name {
		name := model.UserName{
			ID:            srv.bundle.idr.ULID(),
			UserId:        user.ID,
			FirstName:     input.Name.FirstName,
			LastName:      input.Name.LastName,
			PreferredName: input.Name.PreferredName,
		}

		if err := connect.DB(ctx).Create(name).Error; nil != err {
			return err
		}
	}

	return nil
}
