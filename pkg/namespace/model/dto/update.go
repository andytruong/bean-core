package dto

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
}
