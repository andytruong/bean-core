package dto

import (
	"bean/pkg/namespace/model"
	"bean/pkg/util"
)

type DomainNameInput struct {
	Verified *bool   `json:"verified"`
	Value    *string `json:"value"`
	IsActive *bool   `json:"isActive"`
}

type DomainNamesInput struct {
	Primary   *DomainNameInput   `json:"primary"`
	Secondary []*DomainNameInput `json:"secondary"`
}

type NamespaceCreateContext struct {
	UserID string `json:"userId"`
}

type NamespaceCreateInput struct {
	Object  *NamespaceCreateInputObject `json:"object"`
	Context *NamespaceCreateContext     `json:"context"`
}

type NamespaceCreateInputObject struct {
	Title       *string           `json:"title"`
	DomainNames *DomainNamesInput `json:"domainNames"`
}

type NamespaceCreateOutcome struct {
	Errors    []*util.Error    `json:"errors"`
	Namespace *model.Namespace `json:"namespace"`
}
