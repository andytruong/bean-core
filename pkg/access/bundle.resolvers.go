package access

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bean/components/claim"
	"bean/components/util"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	space_model "bean/pkg/space/model"
	user_model "bean/pkg/user/model"
)

func (bundle *AccessBundle) newResolves() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"AccessQuery": func(ctx context.Context) (*dto.AccessQuery, error) {
				return &dto.AccessQuery{}, nil
			},
		},
		"Mutation": map[string]interface{}{
			"AccessMutation": func(ctx context.Context) (*dto.AccessMutation, error) {
				return &dto.AccessMutation{}, nil
			},
		},
		"AccessQuery": map[string]interface{}{
			"Session": func(ctx context.Context, _ *dto.AccessQuery) (*dto.AccessSessionQuery, error) {
				return &dto.AccessSessionQuery{}, nil
			},
		},
		"AccessSessionQuery": map[string]interface{}{
			"Load": func(ctx context.Context, _ *dto.AccessSessionQuery, id string) (*model.Session, error) {
				return bundle.sessionService.load(ctx, id)
			},
		},
		"AccessMutation": map[string]interface{}{
			"Session": func(ctx context.Context, _ *dto.AccessMutation) (*dto.AccessSessionMutation, error) {
				return &dto.AccessSessionMutation{}, nil
			},
		},
		"AccessSessionMutation": map[string]interface{}{
			"Create": func(ctx context.Context, _ *dto.AccessSessionMutation, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
				return bundle.sessionService.Create(ctx, in)
			},
			"Archive": func(ctx context.Context, _ *dto.AccessSessionMutation) (*dto.SessionArchiveOutcome, error) {
				claims := claim.ContextToPayload(ctx)
				if nil == claims {
					return nil, util.ErrorAuthRequired
				}

				sess, _ := bundle.sessionService.load(ctx, claims.SessionId())
				if sess != nil {
					out, err := bundle.sessionService.Delete(ctx, sess)

					return out, err
				}

				return &dto.SessionArchiveOutcome{
					Errors: util.NewErrors(util.ErrorCodeInput, []string{"token"}, "session not found"),
					Result: false,
				}, nil
			},
		},
		"Session": map[string]interface{}{
			"User": func(ctx context.Context, obj *model.Session) (*user_model.User, error) {
				return bundle.userBundle.Service.Load(ctx, obj.UserId)
			},
			"Context": func(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
				panic("implement me")
			},
			"Scopes": func(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
				return obj.Scopes, nil
			},
			"Space": func(ctx context.Context, obj *model.Session) (*space_model.Space, error) {
				return bundle.spaceBundle.Service.Load(ctx, obj.SpaceId)
			},
			"Jwt": func(ctx context.Context, session *model.Session, codeVerifier string) (string, error) {
				roles, err := bundle.spaceBundle.MemberService.FindRoles(ctx, session.UserId, session.SpaceId)

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
						ExpiresAt: time.Now().Add(bundle.cnf.Jwt.Timeout).Unix(),
						Subject:   session.UserId,
						Audience:  session.SpaceId,
					},
				}

				return bundle.JwtService.Sign(claims)
			},
		},
	}
}
