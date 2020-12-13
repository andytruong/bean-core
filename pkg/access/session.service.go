package access

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/components/claim"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mSpace "bean/pkg/space/model"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
)

type SessionService struct {
	bundle *AccessBundle
}

func (this *SessionService) Create(tx *gorm.DB, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
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

func (this *SessionService) createUseCredentials(tx *gorm.DB, in *dto.SessionCreateUseCredentialsInput) (*dto.SessionCreateOutcome, error) {
	// load email object, so we have userID
	email := mUser.UserEmail{}
	{
		err := tx.First(&email, "value = ?", in.Email).Error
		if nil != err {
			return nil, errors.New("userBundle not found")
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
					Errors: util.NewErrors(util.ErrorCodeInput, []string{"input.spaceId"}, "invalid password"),
				}, nil
			}
		}
	}

	return this.create(tx, claim.KindCredentials, email.UserId, in.SpaceID, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
}

func (this *SessionService) generateOTLT(tx *gorm.DB, in *dto.SessionCreateGenerateOTLT) (*dto.SessionCreateOutcome, error) {
	return this.create(tx, claim.KindOTLT, in.UserID, in.SpaceID, nil)
}

func (this *SessionService) useOTLT(tx *gorm.DB, in *dto.SessionCreateUseOTLT) (*dto.SessionCreateOutcome, error) {
	oneTimeSession, err := this.LoadByToken(tx.Statement.Context, tx, in.Token)
	if nil != err {
		return nil, err
	}

	if oneTimeSession.Kind != claim.KindOTLT {
		return nil, util.ErrorInvalidArgument
	}

	out, err := this.create(tx, claim.KindAuthenticated, oneTimeSession.UserId, oneTimeSession.SpaceId, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
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

func (this SessionService) create(
	tx *gorm.DB,
	kind claim.Kind, userId string, spaceId string,
	create func(*model.Session),
) (*dto.SessionCreateOutcome, error) {
	membership := &mSpace.Membership{}

	// validate membership
	{
		err := tx.
			First(&membership, "space_id = ? AND user_id = ?", spaceId, userId).
			Error

		if err == gorm.ErrRecordNotFound {
			return &dto.SessionCreateOutcome{
				Errors: util.NewErrors(util.ErrorCodeInput, []string{"input.spaceId"}, "membership not found"),
			}, nil
		}
	}

	token := this.bundle.id.MustULID()
	session := &model.Session{
		ID:          this.bundle.id.MustULID(),
		Version:     this.bundle.id.MustULID(),
		Kind:        kind,
		UserId:      userId,
		SpaceId:     spaceId,
		HashedToken: this.bundle.id.Encode(token),
		Scopes:      nil, // TODO
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(this.bundle.config.SessionTimeout),
	}

	if nil != create {
		create(session)
	}

	if err := tx.Create(&session).Error; nil != err {
		return nil, err
	} else {
		// update membership -> last-time-login
		err := this.bundle.spaceBundle.MemberService.UpdateLastLoginTime(tx, membership)
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

func (this SessionService) load(ctx context.Context, db *gorm.DB, id string) (*model.Session, error) {
	session := &model.Session{}
	err := db.
		WithContext(ctx).
		First(&session, "id = ?", id).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + id)
	}

	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	if !session.IsActive {
		return nil, errors.New("session archived")
	}

	return session, nil
}

func (this SessionService) LoadByToken(ctx context.Context, db *gorm.DB, token string) (*model.Session, error) {
	session := &model.Session{}
	err := db.
		WithContext(ctx).
		First(&session, "hashed_token = ?", this.bundle.id.Encode(token)).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + this.bundle.id.Encode(token))
	}

	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	if !session.IsActive {
		return nil, errors.New("session archived")
	}

	return session, nil
}

func (this SessionService) Delete(tx *gorm.DB, session *model.Session) (*dto.SessionArchiveOutcome, error) {
	session.IsActive = false
	session.Version = this.bundle.id.MustULID()
	session.UpdatedAt = time.Now()
	err := tx.Save(&session).Error
	if nil != err {
		return nil, err
	} else {
		// If session.kind is … also archive parent sessions
		if session.Kind == claim.KindOTLT {

		}

		// If session.kind is KindCredentials/KindAuthenticated also archive child sessions
		if session.Kind == claim.KindCredentials || session.Kind == claim.KindAuthenticated {
			// find & archive all child sessions
		}
	}

	return &dto.SessionArchiveOutcome{Errors: nil, Result: true}, nil
}
