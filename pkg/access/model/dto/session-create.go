package dto

import (
	"bean/pkg/access/model"
	"bean/pkg/util"
)

type SessionCreateInput struct {
	NamespaceID    *string                    `json:"namespaceId"`
	Email          util.EmailAddress          `json:"email"`
	HashedPassword string                     `json:"hashedPassword"`
	Context        *SessionCreateContextInput `json:"context"`
}

type SessionCreateContextInput struct {
	IPAddress  *string           `json:"ipAddress"`
	Country    *string           `json:"country"`
	DeviceType *model.DeviceType `json:"deviceType"`
	DeviceName *string           `json:"deviceName"`
}

type SessionCreateOutcome struct {
	Errors  []*util.Error  `json:"errors"`
	Session *model.Session `json:"session"`
}
