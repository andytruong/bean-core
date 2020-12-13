package access

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bean/components/claim"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	space_model "bean/pkg/space/model"
	user_model "bean/pkg/user/model"
	"bean/pkg/util"
)

func (this *AccessBundle) newResolves() map[string]interface{} {
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
			"AccessSessionQuery": func(ctx context.Context) (*dto.AccessSessionQuery, error) {
				return &dto.AccessSessionQuery{}, nil
			},
		},
		"AccessSessionQuery": map[string]interface{}{
			"Load": func(ctx context.Context, id string) (*model.Session, error) {
				return this.sessionService.load(ctx, this.db, id)
			},
		},
		"AccessMutation": map[string]interface{}{
			"AccessSessionMutation": func(ctx context.Context) (*dto.AccessSessionMutation, error) {
				return &dto.AccessSessionMutation{}, nil
			},
		},
		"AccessSessionMutation": map[string]interface{}{
			"Create": func(ctx context.Context, in *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
				return this.sessionService.Create(this.db, in)
			},
			"Archive": func(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
				claims := claim.ContextToPayload(ctx)
				if nil == claims {
					return nil, util.ErrorAuthRequired
				}

				sess, _ := this.sessionService.load(ctx, this.db.WithContext(ctx), claims.SessionId())
				if sess != nil {
					out, err := this.sessionService.Delete(this.db.WithContext(ctx), sess)

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
				return this.userBundle.Service.Load(this.db.WithContext(ctx), obj.UserId)
			},
			"Context": func(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
				panic("implement me")
			},
			"Scopes": func(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
				return obj.Scopes, nil
			},
			"Space": func(ctx context.Context, obj *model.Session) (*space_model.Space, error) {
				return this.spaceBundle.Load(ctx, obj.SpaceId)
			},
			"Jwt": func(ctx context.Context, session *model.Session, codeVerifier string) (string, error) {
				roles, err := this.spaceBundle.MemberService.FindRoles(ctx, session.UserId, session.SpaceId)

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

				return this.Sign(claims)
			},
		},
	}
}
