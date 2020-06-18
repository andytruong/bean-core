package user

import (
	"gorm.io/gorm"

	"bean/pkg/user/model"
	"bean/pkg/user/model/dto"
)

type CoreName struct {
	bean *UserBean
}

func (this *CoreName) create(tx *gorm.DB, user *model.User, input *dto.UserCreateInput) error {
	if nil != input.Name {
		name := model.UserName{
			ID:            this.bean.id.MustULID(),
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
