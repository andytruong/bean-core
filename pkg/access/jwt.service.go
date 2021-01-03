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

func (srv JwtService) Sign(claims jwt.Claims) (string, error) {
	key, err := srv.bundle.cnf.GetSignKey()
	if nil != err {
		return "", errors.Wrap(util.ErrorConfig, err.Error())
	}

	return jwt.
		NewWithClaims(srv.bundle.cnf.signMethod(), claims).
		SignedString(key)
}
