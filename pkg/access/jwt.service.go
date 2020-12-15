package access

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"bean/components/claim"
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
