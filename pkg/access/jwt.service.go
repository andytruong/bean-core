package access

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bean/components/claim"
	"bean/components/util"
	"bean/pkg/access/model"
)

func newJwtService(bundle *Bundle) (*JwtService, error) {
	privateKey, err := func(path string) (interface{}, error) {
		file, err := ioutil.ReadFile(path)
		if nil != err {
			return nil, err
		}

		block, _ := pem.Decode(file)

		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}(bundle.cnf.Jwt.PrivateKey.String())

	if nil != err {
		return nil, err
	}

	publicKey, err := claim.ParseRsaPublicKeyFromFile(bundle.cnf.Jwt.PublicKey.String())
	if nil != err {
		return nil, err
	}

	return &JwtService{
		bundle:     bundle,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

type JwtService struct {
	bundle     *Bundle
	privateKey interface{}
	publicKey  interface{}
}

func (srv JwtService) Validate(authHeader string) (*claim.Payload, error) {
	chunks := strings.Split(authHeader, " ")
	authHeader = chunks[len(chunks)-1]

	if parts := strings.Split(authHeader, "."); len(parts) == 3 {
		token, err := jwt.ParseWithClaims(
			authHeader,
			&claim.Payload{},
			func(token *jwt.Token) (interface{}, error) {
				return srv.publicKey, nil
			},
		)

		if nil != err {
			return nil, err
		} else {
			return token.Claims.(*claim.Payload), nil
		}
	}

	return nil, ErrInvalidHeader
}

func (srv JwtService) getSignedString(ctx context.Context, session *model.Session, codeVerifier string) (
	string, error,
) {
	roles, err := srv.bundle.spaceBundle.MemberService.FindRoles(ctx, session.UserId, session.SpaceId)

	if nil != err {
		return "", err
	}

	if !session.Verify(codeVerifier) {
		return "", fmt.Errorf("can not verify")
	}

	payload := claim.NewPayload()
	payload.
		SetKind(session.Kind).
		SetSessionId(session.ID).
		SetUserId(session.UserId).
		SetSpaceId(session.SpaceId).
		SetApplication("access"). // TODO: Change to application ID
		SetExpireAt(time.Now().Add(srv.bundle.cnf.Jwt.Timeout).Unix())

	for _, role := range roles {
		payload.AddRole(role.Title)
	}

	return srv.bundle.JwtService.SignedString(payload)
}

func (srv JwtService) SignedString(claims jwt.Claims) (string, error) {
	return jwt.
		NewWithClaims(srv.signMethod(), claims).
		SignedString(srv.privateKey)
}

func (srv JwtService) signMethod() jwt.SigningMethod {
	switch srv.bundle.cnf.Jwt.Algorithm {
	case "RS512":
		return jwt.SigningMethodRS512

	default:
		panic(util.ErrorToBeImplemented)
	}
}
