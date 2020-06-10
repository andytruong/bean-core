package dto

import (
	"bean/pkg/config/model"
	"bean/pkg/util"
)

type VariableCreateInput struct {
	BucketId    string  `json:"bucketId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Value       string  `json:"value"`
}

type VariableUpdateInput struct {
	Id          string  `json:"id"`
	Version     string  `json:"version"`
	Description *string `json:"description"`
	Value       *string `json:"value"`
}

type VariableDeleteInput struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type VariableMutationOutcome struct {
	Errors   []*util.Error         `json:"errors"`
	Variable *model.ConfigVariable `json:"variable"`
}
