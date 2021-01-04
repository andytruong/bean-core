package dto

import (
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/access/model"
)

type SessionCreateInput struct {
	SpaceID             string              `json:"spaceId"`
	Email               scalar.EmailAddress `json:"email"`
	HashedPassword      string              `json:"hashedPassword"`
	CodeChallengeMethod string              `json:"codeChallengeMethod"`
	CodeChallenge       string              `json:"codeChallenge"`
}

type SessionCreateUseCredentialsInput struct {
}

type SessionCreateOTLTSessionInput struct {
	SpaceID string `json:"spaceId"`
	UserID  string `json:"userId"`
}

type SessionExchangeOTLTInput struct {
	Token               string `json:"token"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
	CodeChallenge       string `json:"codeChallenge"`
}

type SessionCreateContextInput struct {
	IPAddress  *string           `json:"ipAddress"`
	Country    *string           `json:"country"`
	DeviceType *model.DeviceType `json:"deviceType"`
	DeviceName *string           `json:"deviceName"`
}

type SessionOutcome struct {
	Errors  []*util.Error  `json:"errors"`
	Token   *string        `json:"token"`
	Session *model.Session `json:"session"`
}
