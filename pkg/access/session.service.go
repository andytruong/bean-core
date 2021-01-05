package access

import (
	"context"
	"errors"
	"time"
	
	"gorm.io/gorm"
	
	"bean/components/claim"
	"bean/components/connect"
	"bean/components/util"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	mSpace "bean/pkg/space/model"
	"bean/pkg/user"
)

type SessionService struct {
	bundle *Bundle
}

func (srv *SessionService) newSessionWithCredentials(ctx context.Context, in *dto.SessionCreateInput) (
	*dto.SessionOutcome, error,
) {
	// load email object, so we have userID
	email, err := srv.bundle.userBundle.EmailService.Load(ctx, in.Email)
	if nil != err {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrorUserNotFound
		}
		
		return nil, err
	} else if !email.IsActive {
		err := user.ErrorEmailInactive.Error()
		errList := util.NewErrors(util.ErrorCodeInput, []string{"input.email"}, err)
		
		return &dto.SessionOutcome{Errors: errList}, nil
	}
	
	// password validation
	passwordMatched, err := srv.bundle.userBundle.PasswordService.ValidPassword(ctx, email.UserId, in.HashedPassword)
	if nil != err {
		return nil, err
	} else if !passwordMatched {
		errList := util.NewErrors(util.ErrorCodeInput, []string{"input.spaceId"}, "invalid password")
		
		return &dto.SessionOutcome{Errors: errList}, nil
	}
	
	return srv.create(ctx, claim.KindCredentials, email.UserId, in.SpaceID, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
}

func (srv *SessionService) newOTLTSession(ctx context.Context, in *dto.SessionCreateOTLTSessionInput) (
	*dto.SessionOutcome, error,
) {
	return srv.create(ctx, claim.KindOTLT, in.UserID, in.SpaceID, nil)
}

func (srv *SessionService) newSessionWithOTLT(ctx context.Context, in *dto.SessionExchangeOTLTInput) (
	*dto.SessionOutcome, error,
) {
	oneTimeSession, err := srv.LoadByToken(ctx, in.Token)
	if nil != err {
		return nil, err
	}
	
	if oneTimeSession.Kind != claim.KindOTLT {
		return nil, util.ErrorInvalidArgument
	}
	
	out, err := srv.create(ctx, claim.KindAuthenticated, oneTimeSession.UserId, oneTimeSession.SpaceId, func(session *model.Session) {
		session.CodeChallengeMethod = in.CodeChallengeMethod
		session.CodeChallenge = in.CodeChallenge
	})
	if nil != err {
		return nil, err
	}
	
	// delete OTLT session
	{
		_, err := srv.Delete(ctx, oneTimeSession)
		if nil != err {
			return nil, err
		}
	}
	
	return out, err
}

func (srv SessionService) create(
	ctx context.Context,
	kind claim.Kind, userId string, spaceId string,
	create func(*model.Session),
) (*dto.SessionOutcome, error) {
	tx := connect.ContextToDB(ctx)
	membership := &mSpace.Membership{}
	
	// validate membership
	{
		err := tx.
			First(&membership, "space_id = ? AND user_id = ?", spaceId, userId).
			Error
		
		if err == gorm.ErrRecordNotFound {
			return &dto.SessionOutcome{
				Errors: util.NewErrors(util.ErrorCodeInput, []string{"input.spaceId"}, "membership not found"),
			}, nil
		}
	}
	
	token := srv.bundle.idr.ULID()
	session := &model.Session{
		ID:          srv.bundle.idr.ULID(),
		Version:     srv.bundle.idr.ULID(),
		Kind:        kind,
		UserId:      userId,
		SpaceId:     spaceId,
		HashedToken: srv.bundle.idr.Encode(token),
		Scopes:      nil, // TODO
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiredAt:   time.Now().Add(srv.bundle.cnf.SessionTimeout),
	}
	
	if nil != create {
		create(session)
	}
	
	if err := tx.Create(&session).Error; nil != err {
		return nil, err
	} else {
		// update membership -> last-time-login
		err := srv.bundle.spaceBundle.MemberService.UpdateLastLoginTime(ctx, membership)
		if nil != err {
			return nil, err
		}
	}
	
	return &dto.SessionOutcome{
		Errors:  nil,
		Token:   &token,
		Session: session,
	}, nil
}

func (srv SessionService) load(ctx context.Context, id string) (*model.Session, error) {
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

func (srv SessionService) LoadByToken(ctx context.Context, token string) (*model.Session, error) {
	db := connect.ContextToDB(ctx)
	session := &model.Session{}
	err := db.First(&session, "hashed_token = ?", srv.bundle.idr.Encode(token)).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("session not found: " + srv.bundle.idr.Encode(token))
	}
	
	if session.ExpiredAt.Unix() <= time.Now().Unix() {
		return nil, errors.New("session expired")
	}
	
	if !session.IsActive {
		return nil, errors.New("session archived")
	}
	
	return session, nil
}

func (srv SessionService) Delete(ctx context.Context, session *model.Session) (*dto.SessionArchiveOutcome, error) {
	tx := connect.ContextToDB(ctx)
	session.IsActive = false
	session.Version = srv.bundle.idr.ULID()
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
