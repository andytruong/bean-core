package access

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"bean/components/claim"
	"bean/components/util"
)

type JwtService struct {
	bundle *AccessBundle
}

func (service JwtService) Validate(authHeader string) (*claim.Payload, error) {
	chunks := strings.Split(authHeader, " ")
	authHeader = chunks[len(chunks)-1]

	if parts := strings.Split(authHeader, "."); len(parts) == 3 {
		token, err := jwt.ParseWithClaims(
			authHeader,
			&claim.Payload{},
			func(token *jwt.Token) (interface{}, error) {
				return service.bundle.config.GetParseKey()
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

func (service JwtService) Sign(claims jwt.Claims) (string, error) {
	key, err := service.bundle.config.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}

	return jwt.
		NewWithClaims(service.bundle.config.signMethod(), claims).
		SignedString(key)
}
