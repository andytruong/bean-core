package dto

import (
	"bean/pkg/access/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type SessionCreateInput struct {
	UseCredentials *SessionCreateUseCredentialsInput `json:"useCredentials"`
	GenerateOTLT   *SessionCreateGenerateOTLT        `json:"generateOTLT"`
	UseOTLT        *SessionCreateUseOTLT             `json:"useOTLT"`
	Context        *SessionCreateContextInput        `json:"context"`
}

type SessionCreateUseCredentialsInput struct {
	NamespaceID    string           `json:"namespaceId"`
	Email          api.EmailAddress `json:"email"`
	HashedPassword string           `json:"hashedPassword"`
}

type SessionCreateGenerateOTLT struct {
	NamespaceID string `json:"namespaceId"`
	UserID      string `json:"userId"`
}

type SessionCreateUseOTLT struct {
	Token string `json:"token"`
}

type SessionCreateContextInput struct {
	IPAddress  *string           `json:"ipAddress"`
	Country    *string           `json:"country"`
	DeviceType *model.DeviceType `json:"deviceType"`
	DeviceName *string           `json:"deviceName"`
}

type SessionCreateOutcome struct {
	Errors  []*util.Error  `json:"errors"`
	Token   *string        `json:"token"`
	Session *model.Session `json:"session"`
}
