package model

import (
	"time"
)

type Session struct {
	ID          string          `json:"id"`
	HashedToken string          `json:"hashedToken"`
	Scopes      []*AccessScope  `json:"scopes"`
	Context     *SessionContext `json:"context"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	ExpiredAt   time.Time       `json:"expiredAt"`
}

type SessionContext struct {
	IPAddress  *string     `json:"ipAddress"`
	Country    *string     `json:"country"`
	DeviceType *DeviceType `json:"deviceType"`
	DeviceName *string     `json:"deviceName"`
}
