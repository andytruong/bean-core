package user

import (
	"gorm.io/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type NameService struct {
	bundle *UserBundle
}

func (this *NameService) load(db *gorm.DB, userId string) (*model.UserName, error) {
	name := &model.UserName{}
	err := db.Where(model.UserName{UserId: userId}).First(&name).Error
	if nil != err {
		return nil, err
	}

	return name, nil
}

func (this *NameService) create(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Name {
		name := model.UserName{
			ID:            this.bundle.id.MustULID(),
			UserId:        user.ID,
			FirstName:     input.Name.FirstName,
			LastName:      input.Name.LastName,
			PreferredName: input.Name.PreferredName,
		}

		if err := tx.Create(name).Error; nil != err {
			return err
		}
	}

	return nil
}
