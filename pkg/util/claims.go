package util

import (
	"crypto/rsa"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

type (
	// maxLength: 32
	Kind string

	Claims struct {
		jwt.StandardClaims
		Kind  Kind     `json:"kind"`
		Roles []string `json:"roles"`
	}
)

const (
	// session created with user/password
	// with this session, use can access almost endpoints provided for them.
	KindCredentials Kind = "credentials"

	// For SSO login we need generate a kind of session, from where user can obtain a full authenticated session.
	// Example flow:
	//   1. User access login page
	//   2. User click auth with Google /auth/with/google
	//   3. User auth using Google login process.
	//   4. User returned /auth/done/google?code=codeFromGoogle
	//   5. Our server:
	//        - Our server load user information from Google.
	//        - If user record found -> create 'oneTime' session
	//   6. User is redirected to /auth/oneTime/$oneTimeSession.token
	//   7. Our server will generate full authenticated session for user, one-time session is deleted.
	KindOTLT Kind = "onetime"

	// User who simply authenticated but without providing credentials.
	// With this kind of session, user can not change password.
	KindAuthenticated Kind = "authenticated"

	// When user forgets password & request for new one, system create a one-time token and send to their
	// email inbox. From there, they can use that token to generate a new password.
	KindPasswordForgot Kind = "password-forget"

	// With this kind of session, user can only reset their password.
	KindPasswordReset Kind = "password-reset"
)

func (this Claims) UserId() string {
	return this.Subject
}

func (this Claims) SessionId() string {
	return this.Id
}

func (this Claims) NamespaceId() string {
	return this.Audience
}

func (this *Claims) UnmarshalGQL(v interface{}) error {
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

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			fmt.Printf("%v %v", claims.Subject, claims.StandardClaims.ExpiresAt)
		} else {
			return err
		}
	}

	return nil
}

func (this Claims) MarshalGQL(w io.Writer) {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		Foo string `json:"foo"`
		jwt.StandardClaims
	}

	// Create the Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, this)
	ss, err := token.SignedString(mySigningKey)
	fmt.Fprintf(w, "%v %v", ss, err)
}

func ParseRsaPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	content, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}

	block, _ := pem.Decode(content)
	key := &rsa.PublicKey{}
	_, err = asn1.Unmarshal(block.Bytes, key)
	if nil != err {
		return nil, err
	}

	return key, nil
}
