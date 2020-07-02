package fixtures

import (
	"bean/components/scalar"
	"bean/pkg/access/model/dto"
)

func SessionCreateInputFixtureUseCredentials(namespaceId string, email string, hashedPassword string) *dto.SessionCreateInput {
	return &dto.SessionCreateInput{
		UseCredentials: &dto.SessionCreateUseCredentialsInput{
			NamespaceID:    namespaceId,
			Email:          scalar.EmailAddress(email),
			HashedPassword: hashedPassword,
		},
		Context: &dto.SessionCreateContextInput{
			IPAddress:  nil,
			Country:    nil,
			DeviceType: nil,
			DeviceName: nil,
		},
	}
}
