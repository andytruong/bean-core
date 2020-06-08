package handler

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type SessionDeleteHandler struct {
	DB *gorm.DB
	ID *util.Identifier
}

func (this SessionDeleteHandler) Handle(ctx context.Context, session *model.Session) (*dto.SessionDeleteOutcome, error) {
	session.IsActive = false
	session.Version = this.ID.MustULID()
	session.UpdatedAt = time.Now()
	err := this.DB.Table(connect.TableAccessSession).Save(&session).Error
	if nil != err {
		return nil, err
	}

	return &dto.SessionDeleteOutcome{Errors: nil, Result: true}, nil
}
