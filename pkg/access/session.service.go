package access

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/util"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mSpace "bean/pkg/space/model"
	mUser "bean/pkg/user/model"
)

type SessionService struct {
	bundle *AccessBundle
}

func (service *SessionService) Create(tx *gorm.DB, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	if nil != in.UseCredentials {
		return service.createUseCredentials(tx, in.UseCredentials)
	}

	if nil != in.GenerateOTLT {
		return service.generateOTLT(tx, in.GenerateOTLT)
	}

	if nil != in.UseOTLT {
		return service.useOTLT(tx, in.UseOTLT)
	}

	return nil, nil
}

func (service *SessionService) createUseCredentials(tx *gorm.DB, in *dto.SessionCreateUseCredentialsInput) (*dto.SessionCreateOutcome, error) {
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

	return service.create(tx, claim.KindCredentials, email.UserId, in.SpaceID, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
}

func (service *SessionService) generateOTLT(tx *gorm.DB, in *dto.SessionCreateGenerateOTLT) (*dto.SessionCreateOutcome, error) {
	return service.create(tx, claim.KindOTLT, in.UserID, in.SpaceID, nil)
}

func (service *SessionService) useOTLT(tx *gorm.DB, in *dto.SessionCreateUseOTLT) (*dto.SessionCreateOutcome, error) {
	oneTimeSession, err := service.LoadByToken(tx, in.Token)
	if nil != err {
		return nil, err
	}

	if oneTimeSession.Kind != claim.KindOTLT {
		return nil, util.ErrorInvalidArgument
	}

	out, err := service.create(tx, claim.KindAuthenticated, oneTimeSession.UserId, oneTimeSession.SpaceId, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
	if nil != err {
		return nil, err
	}

	// delete OTLT session
	{
		_, err := service.Delete(tx, oneTimeSession)
		if nil != err {
			return nil, err
		}
	}

	return out, err
}

func (service SessionService) create(
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

	token := service.bundle.id.MustULID()
	session := &model.Session{
		ID:          service.bundle.id.MustULID(),
		Version:     service.bundle.id.MustULID(),
		Kind:        kind,
		UserId:      userId,
		SpaceId:     spaceId,
		HashedToken: service.bundle.id.Encode(token),
		Scopes:      nil, // TODO
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(service.bundle.config.SessionTimeout),
	}

	if nil != create {
		create(session)
	}

	if err := tx.Create(&session).Error; nil != err {
		return nil, err
	} else {
		// update membership -> last-time-login
		err := service.bundle.spaceBundle.MemberService.UpdateLastLoginTime(tx, membership)
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

func (service SessionService) load(ctx context.Context, db *gorm.DB, id string) (*model.Session, error) {
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

func (service SessionService) LoadByToken(db *gorm.DB, token string) (*model.Session, error) {
	session := &model.Session{}
	err := db.
		First(&session, "hashed_token = ?", service.bundle.id.Encode(token)).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + service.bundle.id.Encode(token))
	}

	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	if !session.IsActive {
		return nil, errors.New("session archived")
	}

	return session, nil
}

func (service SessionService) Delete(tx *gorm.DB, session *model.Session) (*dto.SessionArchiveOutcome, error) {
	session.IsActive = false
	session.Version = service.bundle.id.MustULID()
	session.UpdatedAt = time.Now()
	err := tx.Save(&session).Error
	if nil != err {
		return nil, err
	}

	// TODO
	// If session.kind is … also archive parent sessions
	// if session.Kind == claim.KindOTLT {}

	// TODO
	// If session.kind is KindCredentials/KindAuthenticated also archive child sessions
	// find & archive all child sessions
	// if session.Kind == claim.KindCredentials || session.Kind == claim.KindAuthenticated {}

	return &dto.SessionArchiveOutcome{Errors: nil, Result: true}, nil
}
