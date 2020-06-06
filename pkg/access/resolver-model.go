package access

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"bean/pkg/access/model"
	namespace_model "bean/pkg/namespace/model"
	user_model "bean/pkg/user/model"
	"bean/pkg/util"
)

type ModelResolver struct {
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
	key, err := this.config.signKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}

	claims := util.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "access",
			Id:        session.ID,
			IssuedAt:  time.Now().UnixNano(),
			ExpiresAt: time.Now().Add(this.config.Jwt.Timeout).UnixNano(),
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
