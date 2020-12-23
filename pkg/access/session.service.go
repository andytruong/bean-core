package access

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/util"
	"bean/components/util/connect"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mSpace "bean/pkg/space/model"
	mUser "bean/pkg/user/model"
)

type SessionService struct {
	bundle *AccessBundle
}

func (service *SessionService) Create(ctx context.Context, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	if nil != in.UseCredentials {
		return service.createUseCredentials(ctx, in.UseCredentials)
	}

	if nil != in.GenerateOTLT {
		return service.generateOTLT(ctx, in.GenerateOTLT)
	}

	if nil != in.UseOTLT {
		return service.useOTLT(ctx, in.UseOTLT)
	}

	return nil, nil
}

func (service *SessionService) createUseCredentials(ctx context.Context, in *dto.SessionCreateUseCredentialsInput) (*dto.SessionCreateOutcome, error) {
	tx := connect.ContextToDB(ctx)

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

	return service.create(ctx, claim.KindCredentials, email.UserId, in.SpaceID, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
}

func (service *SessionService) generateOTLT(ctx context.Context, in *dto.SessionCreateGenerateOTLT) (*dto.SessionCreateOutcome, error) {
	return service.create(ctx, claim.KindOTLT, in.UserID, in.SpaceID, nil)
}

func (service *SessionService) useOTLT(ctx context.Context, in *dto.SessionCreateUseOTLT) (*dto.SessionCreateOutcome, error) {
	oneTimeSession, err := service.LoadByToken(ctx, in.Token)
	if nil != err {
		return nil, err
	}

	if oneTimeSession.Kind != claim.KindOTLT {
		return nil, util.ErrorInvalidArgument
	}

	out, err := service.create(ctx, claim.KindAuthenticated, oneTimeSession.UserId, oneTimeSession.SpaceId, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
	if nil != err {
		return nil, err
	}

	// delete OTLT session
	{
		_, err := service.Delete(ctx, oneTimeSession)
		if nil != err {
			return nil, err
		}
	}

	return out, err
}

func (service SessionService) create(
	ctx context.Context,
	kind claim.Kind, userId string, spaceId string,
	create func(*model.Session),
) (*dto.SessionCreateOutcome, error) {
	tx := connect.ContextToDB(ctx)
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

	token := service.bundle.idr.MustULID()
	session := &model.Session{
		ID:          service.bundle.idr.MustULID(),
		Version:     service.bundle.idr.MustULID(),
		Kind:        kind,
		UserId:      userId,
		SpaceId:     spaceId,
		HashedToken: service.bundle.idr.Encode(token),
		Scopes:      nil, // TODO
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(service.bundle.cnf.SessionTimeout),
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

func (service SessionService) load(ctx context.Context, id string) (*model.Session, error) {
	session := &model.Session{}
	db := connect.ContextToDB(ctx)
	err := db.WithContext(ctx).First(&session, "id = ?", id).Error

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

func (service SessionService) LoadByToken(ctx context.Context, token string) (*model.Session, error) {
	db := connect.ContextToDB(ctx)
	session := &model.Session{}
	err := db.First(&session, "hashed_token = ?", service.bundle.idr.Encode(token)).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + service.bundle.idr.Encode(token))
	}

	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}

	if !session.IsActive {
		return nil, errors.New("session archived")
	}

	return session, nil
}

func (service SessionService) Delete(ctx context.Context, session *model.Session) (*dto.SessionArchiveOutcome, error) {
	tx := connect.ContextToDB(ctx)
	session.IsActive = false
	session.Version = service.bundle.idr.MustULID()
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
