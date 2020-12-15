package access

import (
	"context"
	"fmt"
	"time"
	
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	
	"bean/components/claim"
	"bean/pkg/access/model"
	space_model "bean/pkg/space/model"
	user_model "bean/pkg/user/model"
)

type ModelResolver struct {
	logger *zap.Logger
	bundle *AccessBundle
	config *AccessConfiguration
}

// TODO: Remove this
func (this ModelResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	return this.bundle.userBundle.Resolvers.Query.User(ctx, obj.UserId)
}

// TODO: Remove this
func (this ModelResolver) Context(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
	panic("implement me")
}

// TODO: Remove this
func (this ModelResolver) Scopes(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
	return obj.Scopes, nil
}

// TODO: Remove this
func (this ModelResolver) Space(ctx context.Context, obj *model.Session) (*space_model.Space, error) {
	return this.bundle.spaceBundle.Load(ctx, obj.SpaceId)
}

// TODO: Remove this
func (this ModelResolver) Jwt(ctx context.Context, session *model.Session, codeVerifier string) (string, error) {
	roles, err := this.bundle.spaceBundle.MembershipResolver().FindRoles(ctx, session.UserId, session.SpaceId)
	if nil != err {
		return "", err
	}
	
	if !session.Verify(codeVerifier) {
		return "", fmt.Errorf("can not verify")
	}
	
	claims := claim.Payload{
		Kind: session.Kind,
		Roles: func() []string {
			var roleNames []string
			
			for _, role := range roles {
				roleNames = append(roleNames, role.Title)
			}
			
			return roleNames
		}(),
		StandardClaims: jwt.StandardClaims{
			Issuer:    "access",
			Id:        session.ID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(this.config.Jwt.Timeout).Unix(),
			Subject:   session.UserId,
			Audience:  session.SpaceId,
		},
	}
	
	return this.bundle.Sign(claims)
}
