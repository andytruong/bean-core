package claim

import (
	"context"
	"fmt"
	"io"

	"github.com/dgrijalva/jwt-go"
)

func ContextToPayload(ctx context.Context) *Payload {
	if claims, ok := ctx.Value("bean.claims").(*Payload); ok {
		return claims
	}

	return nil
}

type Payload struct {
	jwt.StandardClaims
	Kind  Kind     `json:"kind"`
	Roles []string `json:"roles"`
}

func (payload Payload) UserId() string {
	return payload.Subject
}

func (payload Payload) SessionId() string {
	return payload.Id
}

func (payload Payload) SpaceId() string {
	return payload.Audience
}

func (payload *Payload) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("JWT must be strings")
	} else {
		token, err := jwt.ParseWithClaims(
			in,
			payload,
			func(token *jwt.Token) (interface{}, error) {
				return []byte("AllYourBase"), nil
			},
		)

		if claims, ok := token.Claims.(*Payload); ok && token.Valid {
			fmt.Printf("%v %v", claims.Subject, claims.StandardClaims.ExpiresAt)
		} else {
			return err
		}
	}

	return nil
}

func (payload Payload) MarshalGQL(w io.Writer) {
	mySigningKey := []byte("AllYourBase")
	
	// Create the Payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString(mySigningKey)
	fmt.Fprintf(w, "%v %v", ss, err)
}
