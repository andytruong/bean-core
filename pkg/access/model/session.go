package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Session struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	ParentId    string    `json:"parentId"`
	Kind        Kind      `json:"kind"`
	UserId      string    `json:"userId"`
	NamespaceId string    `json:"namespaceId"`
	HashedToken string    `json:"hashedToken"`
	Scopes      ScopeList `json:"scopes",sql:"type:text"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ExpiredAt   time.Time `json:"expiredAt"`
}

// maxLength: 32
type Kind string

const (
	// session created with user/password
	// with this session, use can access almost endpoints provided for them.
	KindCredentials Kind = "credentials"

	// For SSO login we need generate a kind of session, from where user can obtain a full authenticated session.
	// Example flow:
	//   1. User access login page
	//   2. User click auth with Google /auth/with/google
	//   3. User auth using Google login process.
	//   4. User returned /auth/done/google?code=codeFromGoogle
	//   5. Our server:
	//        - Our server load user information from Google.
	//        - If user record found -> create 'oneTime' session
	//   6. User is redirected to /auth/oneTime/$oneTimeSession.token
	//   7. Our server will generate full authenticated session for user, one-time session is deleted.
	KindOTLT Kind = "onetime"

	// User who simply authenticated but without providing credentials.
	// With this kind of session, user can not change password.
	KindAuthenticated Kind = "authenticated"

	// When user forgets password & request for new one, system create a one-time token and send to their
	// email inbox. From there, they can use that token to generate a new password.
	KindPasswordForgot Kind = "password-forget"

	// With this kind of session, user can only reset their password.
	KindPasswordReset Kind = "password-reset"
)

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
