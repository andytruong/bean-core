package dto

import (
	"bean/pkg/access/model"
	user_model "bean/pkg/user/model"
	"bean/pkg/util"
)

type ValidationOutcome struct {
	Status bool     `json:"status"`
	Errors []*util.Error `json:"errors"`
}

type UserCreateOutcome struct {
	User   *user_model.User `json:"user"`
	Errors []*util.Error    `json:"errors"`
}

type LoginContextInput struct {
	IPAddress  *string           `json:"ipAddress"`
	Country    *string           `json:"country"`
	DeviceType *model.DeviceType `json:"deviceType"`
	DeviceName *string           `json:"deviceName"`
}

type LoginInput struct {
	NamespaceID    *string            `json:"namespaceId"`
	Username       string             `json:"username"`
	HashedPassword string             `json:"hashedPassword"`
	Context        *LoginContextInput `json:"context"`
}

type LogoutInput struct {
	HashedToken string `json:"hashedToken"`
}

type LogoutPayload struct {
	Errors []*util.Error `json:"errors"`
	Result *bool    `json:"result"`
}
