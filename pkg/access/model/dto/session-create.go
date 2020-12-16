package dto

import (
	"bean/components/scalar"
	util2 "bean/components/util"
	"bean/pkg/access/model"
)

type SessionCreateInput struct {
	UseCredentials *SessionCreateUseCredentialsInput `json:"useCredentials"`
	GenerateOTLT   *SessionCreateGenerateOTLT        `json:"generateOTLT"`
	UseOTLT        *SessionCreateUseOTLT             `json:"useOTLT"`
	Context        *SessionCreateContextInput        `json:"context"`
}

type SessionCreateUseCredentialsInput struct {
	SpaceID             string              `json:"spaceId"`
	Email               scalar.EmailAddress `json:"email"`
	HashedPassword      string              `json:"hashedPassword"`
	CodeChallengeMethod string              `json:"codeChallengeMethod"`
	CodeChallenge       string              `json:"codeChallenge"`
}

type SessionCreateGenerateOTLT struct {
	SpaceID string `json:"spaceId"`
	UserID  string `json:"userId"`
}

type SessionCreateUseOTLT struct {
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

type SessionCreateOutcome struct {
	Errors  []*util2.Error `json:"errors"`
	Token   *string        `json:"token"`
	Session *model.Session `json:"session"`
}
