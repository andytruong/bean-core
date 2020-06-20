package access

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mNamespace "bean/pkg/namespace/model"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type Core struct {
	bean *AccessBean
}

func (this Core) Create(tx *gorm.DB, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	// load email object, so we have userID
	email := mUser.UserEmail{}
	{
		err := tx.First(&email, "value = ?", in.Email).Error
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
	membership := &mNamespace.Membership{}

	// password validation
	eg.Go(func() error {
		pass := mUser.UserPassword{}
		err := tx.First(&pass, "user_id = ? AND hashed_value = ? AND is_active = ?", email.UserId, in.HashedPassword, true).Error
		if err == gorm.ErrRecordNotFound {
			outcome.Errors = util.NewErrors(util.ErrorCodeInput, []string{"input.namespaceId"}, "invalid password")
			return nil
		}

		return err
	})

	// membership validation
	eg.Go(func() error {
		err := tx.
			Table(connect.TableNamespaceMemberships).
			First(&membership, "namespace_id = ? AND user_id = ?", in.NamespaceID, email.UserId).
			Error

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

	if nil != outcome.Errors {
		return outcome, nil
	}

	return this.createSession(tx, email.UserId, in.NamespaceID, membership)
}

func (this Core) createSession(
	tx *gorm.DB,
	userId string,
	namespaceId string,
	membership *mNamespace.Membership,
) (*dto.SessionCreateOutcome, error) {
	token := this.bean.id.MustULID()
	session := &model.Session{
		ID:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		UserId:      userId,
		NamespaceId: namespaceId,
		HashedToken: this.bean.id.Encode(token),
		Scopes:      nil, // TODO
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(this.bean.config.SessionTimeout),
	}

	if err := tx.Table(connect.TableAccessSession).Create(&session).Error; nil != err {
		return nil, err
	} else {
		// update membership -> last-time-login
		err := this.bean.namespace.MembershipResolver().UpdateLastLoginTime(tx, membership)
		if nil != err {
			return nil, err
		}
	}

	return &dto.SessionCreateOutcome{
		Errors:  nil,
		Token:   &token,
		Session: session,
	}, nil
}

func (this Core) Load(ctx context.Context, token string) (*model.Session, error) {
	session := &model.Session{}
	err := this.bean.db.
		WithContext(ctx).
		Table(connect.TableAccessSession).
		First(&session, "hashed_token = ?", this.bean.id.Encode(token)).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + this.bean.id.Encode(token))
	}

	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	if !session.IsActive {
		return nil, errors.New("session archived")
	}

	return session, nil
}

func (this Core) Delete(ctx context.Context, session *model.Session) (*dto.SessionDeleteOutcome, error) {
	session.IsActive = false
	session.Version = this.bean.id.MustULID()
	session.UpdatedAt = time.Now()
	err := this.bean.db.WithContext(ctx).Table(connect.TableAccessSession).Save(&session).Error
	if nil != err {
		return nil, err
	}

	return &dto.SessionDeleteOutcome{Errors: nil, Result: true}, nil
}
