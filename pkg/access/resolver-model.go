package access

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"bean/pkg/access/model"
	namespace_model "bean/pkg/namespace/model"
	user_model "bean/pkg/user/model"
	"bean/pkg/util"
)

type ModelResolver struct {
	logger *zap.Logger
	module *AccessModule
	config *Config
}

func (this ModelResolver) User(ctx context.Context, obj *model.Session) (*user_model.User, error) {
	return this.module.userModule.User(ctx, obj.UserId)
}

func (this ModelResolver) Context(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
	panic("implement me")
}

func (this ModelResolver) Scopes(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
	return obj.Scopes, nil
}

func (this ModelResolver) Namespace(ctx context.Context, obj *model.Session) (*namespace_model.Namespace, error) {
	return this.module.namespaceModule.Namespace(ctx, obj.NamespaceId)
}

func (this ModelResolver) Jwt(ctx context.Context, session *model.Session) (string, error) {
	key, err := this.config.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}

	claims := util.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "access",
			Id:        session.ID,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(this.config.Jwt.Timeout).Unix(),
			Subject:   session.UserId,
			Audience:  session.NamespaceId,
		},
	}

	signedString, err := jwt.
		NewWithClaims(this.config.signMethod(), claims).
		SignedString(key)

	if nil != err {
		return "", err
	}

	return signedString, err
}

func (this ModelResolver) JwtValidation(authHeader string) (*util.Claims, error) {
	chunks := strings.Split(authHeader, " ")
	authHeader = chunks[len(chunks)-1]

	if parts := strings.Split(authHeader, "."); 3 == len(parts) {
		token, err := jwt.ParseWithClaims(
			authHeader,
			&util.Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return this.config.GetParseKey()
			},
		)

		if nil != err {
			return nil, err
		} else {
			return token.Claims.(*util.Claims), nil
		}
	}

	return nil, fmt.Errorf("ivnalid authentication header")
}
