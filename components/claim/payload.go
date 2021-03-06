package claim

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bean/components/scalar"
)

func NewPayload() Payload {
	return Payload{
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
		Roles: []string{},
	}
}

func ContextToPayload(ctx context.Context) *Payload {
	if claims, ok := ctx.Value(ClaimsContextKey).(*Payload); ok {
		return claims
	}

	return nil
}

const ClaimsContextKey scalar.ContextKey = "bean.claims"

// Id       -> sessionID
// Issuer   -> applicationID
// Subject  -> userID
// Audience -> spaceID
type Payload struct {
	jwt.StandardClaims
	Kind  Kind     `json:"kind"`
	Roles []string `json:"roles,omitempty"`
}

func (pl *Payload) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, pl)
}

func (pl *Payload) SetKind(value Kind) *Payload {
	pl.Kind = value
	return pl
}

func (pl *Payload) SetExpireAt(value int64) *Payload {
	pl.ExpiresAt = value
	return pl
}

func (pl *Payload) SetApplication(value string) *Payload {
	pl.Issuer = value
	return pl
}

func (pl *Payload) GetApplication() string {
	return pl.Issuer
}

func (pl *Payload) SetUserId(value string) *Payload {
	pl.Subject = value
	return pl
}

func (pl Payload) UserId() string {
	return pl.Subject
}

func (pl *Payload) SetSessionId(value string) *Payload {
	pl.Id = value
	return pl
}

func (pl Payload) SessionId() string {
	return pl.Id
}

func (pl *Payload) SetSpaceId(value string) *Payload {
	pl.Audience = value

	return pl
}

func (pl Payload) SpaceId() string {
	return pl.Audience
}

func (pl *Payload) AddRole(values ...string) *Payload {
	pl.Roles = append(pl.Roles, values...)

	return pl
}

func (pl *Payload) UnmarshalGQL(v interface{}) error {
	if in, ok := v.(string); !ok {
		return fmt.Errorf("JWT must be strings")
	} else {
		token, err := jwt.ParseWithClaims(
			in,
			pl,
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

func (pl Payload) MarshalGQL(w io.Writer) {
	mySigningKey := []byte("AllYourBase")

	// Create the Payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, pl)
	ss, err := token.SignedString(mySigningKey)
	fmt.Fprintf(w, "%v %v", ss, err)
}
