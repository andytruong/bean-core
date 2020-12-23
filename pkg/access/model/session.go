package model

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"bean/components/claim"
)

type Session struct {
	ID                  string     `json:"id"`
	Version             string     `json:"version"`
	ParentId            string     `json:"parentId"`
	Kind                claim.Kind `json:"kind"`
	UserId              string     `json:"userId"`
	SpaceId             string     `json:"spaceId"`
	HashedToken         string     `json:"hashedToken"`
	Scopes              ScopeList  `json:"scopes" sql:"type:text"`
	IsActive            bool       `json:"isActive"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	ExpiredAt           time.Time  `json:"expiredAt"`
	CodeChallengeMethod string     `json:"codeChallengeMethod"`
	CodeChallenge       string     `json:"codeChallenge"`
}

func (Session) TableName() string {
	return "access_session"
}

func (session Session) Verify(codeVerifier string) bool {
	switch session.CodeChallengeMethod {
	case "S256":
		actual := fmt.Sprintf("%x", sha256.Sum256([]byte(codeVerifier)))

		return session.CodeChallenge == actual

	default:
		return false
	}
}

type ScopeList []*AccessScope

func (list ScopeList) Value() (driver.Value, error) {
	bytes, err := json.Marshal(list)
	return string(bytes), err
}

func (list *ScopeList) Scan(in interface{}) error {
	switch value := in.(type) {
	case string:
		return json.Unmarshal([]byte(value), list)
	case []byte:
		return json.Unmarshal(value, list)
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
