package model

import (
	"time"

	"bean/pkg/namespace/model"
	model2 "bean/pkg/user/model"
)

type Session struct {
	HashedToken string           `json:"hashedToken"`
	User        *model2.User     `json:"user"`
	Namespace   *model.Namespace `json:"namespace"`
	Scopes      []*AccessScope   `json:"scopes"`
	Context     *SessionContext  `json:"context"`
	IsActive    bool             `json:"isActive"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	ExpiredAt   time.Time        `json:"expiredAt"`
}

type SessionContext struct {
	IPAddress  *string     `json:"ipAddress"`
	Country    *string     `json:"country"`
	DeviceType *DeviceType `json:"deviceType"`
	DeviceName *string     `json:"deviceName"`
}
