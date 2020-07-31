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

func (this Payload) UserId() string {
	return this.Subject
}

func (this Payload) SessionId() string {
	return this.Id
}

func (this Payload) SpaceId() string {
	return this.Audience
}

func (this *Payload) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("JWT must be strings")
	} else {
		token, err := jwt.ParseWithClaims(
			in,
			this,
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

func (this Payload) MarshalGQL(w io.Writer) {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// Create the Payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, this)
	ss, err := token.SignedString(mySigningKey)
	fmt.Fprintf(w, "%v %v", ss, err)
}
