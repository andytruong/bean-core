package dto

import (
	"bean/pkg/access/model"
	"bean/pkg/util"
)

type SessionCreateInput struct {
	Credentials  *SessionCreateUseCredentialsInput  `json:"credentials"`
	OneTimeLogin *SessionCreateUseOneTimeLoginInput `json:"oneTimeLogin"`
	Context      *SessionCreateContextInput         `json:"context"`
}

type SessionCreateUseCredentialsInput struct {
	NamespaceID    string            `json:"namespaceId"`
	Email          util.EmailAddress `json:"email"`
	HashedPassword string            `json:"hashedPassword"`
}

type SessionCreateUseOneTimeLoginInput struct {
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
