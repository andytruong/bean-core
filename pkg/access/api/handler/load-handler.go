package handler

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/pkg/access/model"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type SessionLoadHandler struct {
	ID *util.Identifier
	DB *gorm.DB
}

func (this SessionLoadHandler) Handle(ctx context.Context, token string) (*model.Session, error) {
	session := &model.Session{}
	err := this.DB.
		WithContext(ctx).
		Table(connect.TableAccessSession).
		First(&session, "hashed_token = ?", this.ID.Encode(token)).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + this.ID.Encode(token))
	}

	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	if !session.IsActive {
		return nil, errors.New("session archived")
	}

	return session, nil
}
