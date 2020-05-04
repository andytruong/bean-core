package dto

import (
	"bean/pkg/access/model"
	"bean/pkg/util"
)

type LoginInput struct {
	NamespaceID    *string            `json:"namespaceId"`
	Username       string             `json:"username"`
	HashedPassword string             `json:"hashedPassword"`
	Context        *LoginContextInput `json:"context"`
}

type LoginContextInput struct {
	IPAddress  *string           `json:"ipAddress"`
	Country    *string           `json:"country"`
	DeviceType *model.DeviceType `json:"deviceType"`
	DeviceName *string           `json:"deviceName"`
}

type LoginOutcome struct {
	Errors  []*util.Error  `json:"errors"`
	Session *model.Session `json:"session"`
}
