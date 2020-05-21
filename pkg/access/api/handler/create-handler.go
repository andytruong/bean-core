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
			return &dto.SessionCreateOutcome{
				Errors: util.NewErrors(util.ErrorCodeInput, []string{"input.email"}, "email address is not active"),
			}, nil
		}
	}

	eg := errgroup.Group{}
	outcome := &dto.SessionCreateOutcome{}

	// password validation
	eg.Go(func() error {
		pass := mUser.UserPassword{}
		err := tx.First(&pass, "user_id = ? AND hashed_value = ? AND is_active = ?", email.UserId, input.Email, true).Error

		if err == gorm.ErrRecordNotFound {
			outcome.Errors = util.NewErrors(util.ErrorCodeInput, []string{"input.namespaceId"}, "invalid password")
			return nil
		}

		return err
	})

	// membership validation
	eg.Go(func() error {
		membership := mNamespace.Membership{}
		err := tx.Table("namespace_memberships").First(&membership, "namespace_id = ? AND user_id = ?", input.NamespaceID, email.UserId).Error

		if err == gorm.ErrRecordNotFound {
			outcome.Errors = util.NewErrors(util.ErrorCodeInput, []string{"input.namespaceId"}, "membership not found")
			return nil
		}

		return err
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
	} else if version, err := this.ID.ULID(); nil != err {
		return nil, err
	} else if token, err := this.ID.UUID(); nil != err {
		return nil, err
	} else {
		session := &model.Session{
			ID:          id,
			Version:     version,
			UserId:      userId,
			NamespaceId: namespaceId,
			HashedToken: this.ID.Encode(token),
			Scopes:      nil, // TODO
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			ExpiredAt:   time.Now(),
		}

		if err := tx.Table("access_session").Create(&session).Error; nil != err {
			return nil, err
		}

		return &dto.SessionCreateOutcome{
			Errors:  nil,
			Token:   &token,
			Session: session,
		}, nil
	}
}
