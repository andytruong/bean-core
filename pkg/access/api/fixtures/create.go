package fixtures

import (
	"crypto/sha256"
	"fmt"

	"bean/components/scalar"
	"bean/pkg/access/model/dto"
)

func SessionCreateInputFixtureUseCredentials(spaceId string, email string, hashedPassword string) *dto.SessionCreateInput {
	codeVerifier := []byte("0123456789")

	return &dto.SessionCreateInput{
		SpaceID:             spaceId,
		Email:               scalar.EmailAddress(email),
		HashedPassword:      hashedPassword,
		CodeChallengeMethod: "S256",
		CodeChallenge:       fmt.Sprintf("%x", sha256.Sum256(codeVerifier)),
	}
}
