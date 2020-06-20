package fixtures

import (
	"bean/pkg/access/model/dto"
	"bean/pkg/util"
)

func SessionCreateInputFixtureUseCredentials(namespaceId string, email string, hashedPassword string) *dto.SessionCreateInput {
	return &dto.SessionCreateInput{
		UseCredentials: &dto.SessionCreateUseCredentialsInput{
			NamespaceID:    namespaceId,
			Email:          util.EmailAddress(email),
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
