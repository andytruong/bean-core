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

type Claims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

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
