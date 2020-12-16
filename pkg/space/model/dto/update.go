package dto

import (
	"bean/pkg/infra/api"
)

type SpaceUpdateInput struct {
	SpaceID      string                  `json:"spaceId"`
	SpaceVersion string                  `json:"spaceVersion"`
	Object       *SpaceUpdateInputObject `json:"object"`
}

type SpaceUpdateInputFeatures struct {
	Register *bool `json:"register"`
}

type SpaceUpdateInputObject struct {
	Features *SpaceUpdateInputFeatures `json:"features"`
	Language *api.Language             `json:"language"`
}
