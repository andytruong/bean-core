package access

import (
	"context"
	"fmt"
	"strings"
	"time"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	
	"bean/components/claim"
	"bean/components/util"
	"bean/pkg/access/model"
)

type JwtService struct {
	bundle *Bundle
}

func (srv JwtService) Validate(authHeader string) (*claim.Payload, error) {
	chunks := strings.Split(authHeader, " ")
	authHeader = chunks[len(chunks)-1]
	
	if parts := strings.Split(authHeader, "."); len(parts) == 3 {
		token, err := jwt.ParseWithClaims(
			authHeader,
			&claim.Payload{},
			func(token *jwt.Token) (interface{}, error) {
				return srv.bundle.cnf.GetParseKey()
			},
		)
		
		if nil != err {
			return nil, err
		} else {
			return token.Claims.(*claim.Payload), nil
		}
	}
	
	return nil, fmt.Errorf("invalid authentication header")
}

func (srv JwtService) claim(ctx context.Context, session *model.Session, codeVerifier string) (string, error) {
	roles, err := srv.bundle.spaceBundle.MemberService.FindRoles(ctx, session.UserId, session.SpaceId)
	
	if nil != err {
		return "", err
	}
	
	if !session.Verify(codeVerifier) {
		return "", fmt.Errorf("can not verify")
	}
	
	claims := claim.NewPayload()
	claims.
		SetKind(session.Kind).
		SetSessionId(session.ID).
		SetUserId(session.UserId).
		SetSpaceId(session.SpaceId).
		SetApplication("access"). // TODO: Change to application ID
		SetExpireAt(time.Now().Add(srv.bundle.cnf.Jwt.Timeout).Unix())
	
	for _, role := range roles {
		claims.AddRole(role.Title)
	}
	
	return srv.bundle.JwtService.Sign(claims)
}

func (srv JwtService) Sign(claims jwt.Claims) (string, error) {
	key, err := srv.bundle.cnf.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}
	
	return jwt.
		NewWithClaims(srv.bundle.cnf.signMethod(), claims).
		SignedString(key)
}
