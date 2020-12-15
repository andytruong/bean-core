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
)

func (this AccessBundle) newResolves() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{},
		"Mutation": map[string]interface{}{
			"SessionCreate": func(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
				return this.SessionCreate(ctx, input)
			},
			"SessionArchive": func(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
				return this.SessionArchive(ctx)
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
