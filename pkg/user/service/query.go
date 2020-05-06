package service

import (
	"github.com/jinzhu/gorm"

	"bean/pkg/user/model"
)

type UserQueryAPI struct {
}

func (this *UserQueryAPI) Load(db *gorm.DB, id string) (*model.User, error) {
	user := &model.User{}
	err := db.Where(&model.User{ID: id}).First(user).Error

	if nil != err {
		return nil, err
	}

	return user, nil
}
