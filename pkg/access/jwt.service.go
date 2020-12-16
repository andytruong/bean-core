package access

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"bean/components/claim"
	util2 "bean/components/util"
)

type JwtService struct {
	bundle *AccessBundle
}

func (this JwtService) Validate(authHeader string) (*claim.Payload, error) {
	chunks := strings.Split(authHeader, " ")
	authHeader = chunks[len(chunks)-1]

	if parts := strings.Split(authHeader, "."); 3 == len(parts) {
		token, err := jwt.ParseWithClaims(
			authHeader,
			&claim.Payload{},
			func(token *jwt.Token) (interface{}, error) {
				return this.bundle.config.GetParseKey()
			},
		)

		if nil != err {
			return nil, err
		} else {
			return token.Claims.(*claim.Payload), nil
		}
	}

	return nil, fmt.Errorf("ivnalid authentication header")
}

func (this JwtService) Sign(claims jwt.Claims) (string, error) {
	key, err := this.bundle.config.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util2.ErrorConfig, err.Error())
	}

	return jwt.
		NewWithClaims(this.bundle.config.signMethod(), claims).
		SignedString(key)
}
