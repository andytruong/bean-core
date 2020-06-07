package dto

import (
	"bean/pkg/namespace/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type NamespaceCreateInput struct {
	Object  NamespaceCreateInputObject `json:"object"`
	Context NamespaceCreateContext     `json:"context"`
}

type NamespaceCreateInputObject struct {
	Kind        model.NamespaceKind    `json:"kind"`
	Title       *string                `json:"title"`
	Language    api.Language           `json:"language"`
	IsActive    bool                   `json:"isActive"`
	DomainNames *DomainNamesInput      `json:"domainNames"`
	Features    NamespaceFeaturesInput `json:"features"`
}

type DomainNameInput struct {
	Verified *bool   `json:"verified"`
	Value    *string `json:"value"`
	IsActive *bool   `json:"isActive"`
}

type DomainNamesInput struct {
	Primary   *DomainNameInput   `json:"primary"`
	Secondary []*DomainNameInput `json:"secondary"`
}

type NamespaceFeaturesInput struct {
	Register bool `json:"register"`
}

type NamespaceCreateContext struct {
	UserID      string  `json:"userId"`
	NamespaceID *string `json:"namespaceId"`
}

type NamespaceCreateOutcome struct {
	Errors    []util.Error     `json:"errors"`
	Namespace *model.Namespace `json:"namespace"`
}
