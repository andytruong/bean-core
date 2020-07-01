package dto

import (
	"bean/pkg/util/api"
)

type NamespaceUpdateInput struct {
	NamespaceID      string                      `json:"namespaceId"`
	NamespaceVersion string                      `json:"namespaceVersion"`
	Object           *NamespaceUpdateInputObject `json:"object"`
}

type NamespaceUpdateInputFeatures struct {
	Register *bool `json:"register"`
}

type NamespaceUpdateInputObject struct {
	Features *NamespaceUpdateInputFeatures `json:"features"`
	Language *api.Language                 `json:"language"`
}
