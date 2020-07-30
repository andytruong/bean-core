package fixtures

import (
	"crypto/sha256"
	"fmt"

	"bean/components/scalar"
	"bean/pkg/access/model/dto"
)

func SessionCreateInputFixtureUseCredentials(namespaceId string, email string, hashedPassword string) *dto.SessionCreateInput {
	return &dto.SessionCreateInput{
		UseCredentials: &dto.SessionCreateUseCredentialsInput{
			NamespaceID:         namespaceId,
			Email:               scalar.EmailAddress(email),
			HashedPassword:      hashedPassword,
			CodeChallengeMethod: "S256",
			CodeChallenge:       fmt.Sprintf("%x", sha256.Sum256([]byte("0123456789"))),
		},
		Context: &dto.SessionCreateContextInput{
			IPAddress:  nil,
			Country:    nil,
			DeviceType: nil,
			DeviceName: nil,
		},
	}
}
