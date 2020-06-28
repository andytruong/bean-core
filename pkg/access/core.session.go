package access

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mNamespace "bean/pkg/namespace/model"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type CoreSession struct {
	bean *AccessBean
}

func (this *CoreSession) Create(tx *gorm.DB, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	if nil != in.UseCredentials {
		return this.createUseCredentials(tx, in.UseCredentials)
	}

	if nil != in.GenerateOTLT {
		return this.generateOTLT(tx, in.GenerateOTLT)
	}

	if nil != in.UseOTLT {
		return this.useOTLT(tx, in.UseOTLT)
	}

	return nil, nil
}

func (this *CoreSession) createUseCredentials(tx *gorm.DB, in *dto.SessionCreateUseCredentialsInput) (*dto.SessionCreateOutcome, error) {
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

	// password validation
	{
		pass := mUser.UserPassword{}
		err := tx.First(&pass, "user_id = ? AND hashed_value = ? AND is_active = ?", email.UserId, in.HashedPassword, true).Error
		if nil != err {
			if err == gorm.ErrRecordNotFound {
				return &dto.SessionCreateOutcome{
					Errors: util.NewErrors(util.ErrorCodeInput, []string{"input.namespaceId"}, "invalid password"),
				}, nil
			}
		}
	}

	return this.create(tx, util.KindCredentials, email.UserId, in.NamespaceID)
}

func (this *CoreSession) generateOTLT(tx *gorm.DB, in *dto.SessionCreateGenerateOTLT) (*dto.SessionCreateOutcome, error) {
	return this.create(tx, util.KindOTLT, in.UserID, in.NamespaceID)
}

func (this *CoreSession) useOTLT(tx *gorm.DB, in *dto.SessionCreateUseOTLT) (*dto.SessionCreateOutcome, error) {
	oneTimeSession, err := this.Load(tx.Statement.Context, tx, in.Token)
	if nil != err {
		return nil, err
	}

	if oneTimeSession.Kind != util.KindOTLT {
		return nil, util.ErrorInvalidArgument
	}

	out, err := this.create(tx, util.KindAuthenticated, oneTimeSession.UserId, oneTimeSession.NamespaceId)
	if nil != err {
		return nil, err
	}

	// delete OTLT session
	{
		_, err := this.Delete(tx, oneTimeSession)
		if nil != err {
			return nil, err
		}
	}

	return out, err
}

func (this CoreSession) create(
	tx *gorm.DB,
	kind util.Kind,
	userId string,
	namespaceId string,
) (*dto.SessionCreateOutcome, error) {
	membership := &mNamespace.Membership{}

	// validate membership
	{
		err := tx.
			Table(connect.TableNamespaceMemberships).
			First(&membership, "namespace_id = ? AND user_id = ?", namespaceId, userId).
			Error

		if err == gorm.ErrRecordNotFound {
			return &dto.SessionCreateOutcome{
				Errors: util.NewErrors(util.ErrorCodeInput, []string{"input.namespaceId"}, "membership not found"),
			}, nil
		}
	}

	token := this.bean.id.MustULID()
	session := &model.Session{
		ID:          this.bean.id.MustULID(),
		Version:     this.bean.id.MustULID(),
		Kind:        kind,
		UserId:      userId,
		NamespaceId: namespaceId,
		HashedToken: this.bean.id.Encode(token),
		Scopes:      nil, // TODO
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(this.bean.genetic.SessionTimeout),
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

func (this CoreSession) Load(ctx context.Context, db *gorm.DB, token string) (*model.Session, error) {
	session := &model.Session{}
	err := db.
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

func (this CoreSession) Delete(tx *gorm.DB, session *model.Session) (*dto.SessionDeleteOutcome, error) {
	session.IsActive = false
	session.Version = this.bean.id.MustULID()
	session.UpdatedAt = time.Now()
	err := tx.Table(connect.TableAccessSession).Save(&session).Error
	if nil != err {
		return nil, err
	}

	return &dto.SessionDeleteOutcome{Errors: nil, Result: true}, nil
}
