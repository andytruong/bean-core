package dto

import (
	"bean/pkg/space/model"
	"bean/pkg/util"
	"bean/pkg/util/api"
)

type SpaceCreateInput struct {
	Object SpaceCreateInputObject `json:"object"`
}

type SpaceCreateInputObject struct {
	Kind        model.SpaceKind    `json:"kind"`
	Title       *string            `json:"title"`
	Language    api.Language       `json:"language"`
	IsActive    bool               `json:"isActive"`
	DomainNames *DomainNamesInput  `json:"domainNames"`
	Features    SpaceFeaturesInput `json:"features"`

	// Internal field
	ParentId *string `json:"parentId"`
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

type SpaceFeaturesInput struct {
	Register bool `json:"register"`
}

type SpaceCreateOutcome struct {
	Errors []util.Error `json:"errors"`
	Space  *model.Space `json:"space"`
}
