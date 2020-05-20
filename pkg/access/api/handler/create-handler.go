package handler

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/sync/errgroup"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mNamespace "bean/pkg/namespace/model"
	"bean/pkg/user"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
)

type SessionCreateHandler struct {
	ID         *util.Identifier
	UserModule *user.UserModule
}

func (this SessionCreateHandler) Handle(ctx context.Context, tx *gorm.DB, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	// load email object, so we have userID
	email := mUser.UserEmail{}
	{
		err := tx.First(&email, "value = ?", input.Email).Error
		if nil != err {
			return nil, errors.New("user not found")
		}

		if !email.IsActive {
			return nil, errors.New("invalid user")
		}
	}

	eg := errgroup.Group{}

	// password validation
	eg.Go(func() error {
		pass := mUser.UserPassword{}
		return tx.First(&pass, "value = ? AND is_active = ?", input.Email, true).Error
	})

	// membership validation
	eg.Go(func() error {
		membership := mNamespace.Membership{}
		return tx.First(&membership, "namespace_id = ? AND user_id = ?", input.NamespaceID, email.UserId).Error
	})

	err := eg.Wait()
	if nil != err {
		return nil, err
	}

	return this.createSession(tx, email.UserId, input.NamespaceID)
}

func (this SessionCreateHandler) createSession(tx *gorm.DB, userId string, namespaceId string) (*dto.SessionCreateOutcome, error) {
	if id, err := this.ID.ULID(); nil != err {
		return nil, err
	} else if token, err := this.ID.UUID(); nil != err {
		return nil, err
	} else if hashedToken, err := this.ID.HashHex("access.token", token); nil != err {
		return nil, err
	} else {
		session := &model.Session{
			ID:          id,
			UserId:      userId,
			NamespaceId: namespaceId,
			HashedToken: hashedToken,
			Scopes:      nil, // TODO
			Context:     nil, // TODO
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ExpiredAt:   time.Now(),
		}

		if err := tx.Create(&session).Error; nil != err {
			return nil, err
		}

		return &dto.SessionCreateOutcome{
			Errors:  nil,
			Token:   &token,
			Session: session,
		}, nil
	}
}
