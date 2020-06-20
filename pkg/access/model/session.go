package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"bean/pkg/util"
)

type Session struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	ParentId    string    `json:"parentId"`
	Kind        util.Kind `json:"kind"`
	UserId      string    `json:"userId"`
	NamespaceId string    `json:"namespaceId"`
	HashedToken string    `json:"hashedToken"`
	Scopes      ScopeList `json:"scopes",sql:"type:text"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ExpiredAt   time.Time `json:"expiredAt"`
}

type ScopeList []*AccessScope

func (this ScopeList) Value() (driver.Value, error) {
	bytes, err := json.Marshal(this)
	return string(bytes), err
}

func (this *ScopeList) Scan(in interface{}) error {
	switch value := in.(type) {
	case string:
		return json.Unmarshal([]byte(value), this)
	case []byte:
		return json.Unmarshal(value, this)
	default:
		return errors.New("not supported")
	}
}

type SessionContext struct {
	IPAddress  *string     `json:"ipAddress"`
	Country    *string     `json:"country"`
	DeviceType *DeviceType `json:"deviceType"`
	DeviceName *string     `json:"deviceName"`
}
