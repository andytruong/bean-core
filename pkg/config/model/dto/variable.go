package dto

import (
	util2 "bean/components/util"
	"bean/pkg/config/model"
)

type VariableCreateInput struct {
	BucketId    string  `json:"bucketId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Value       string  `json:"value"`
	IsLocked    *bool   `json:"isLocked"`
}

type VariableUpdateInput struct {
	Id          string  `json:"id"`
	Version     string  `json:"version"`
	Description *string `json:"description"`
	Value       *string `json:"value"`
	IsLocked    *bool   `json:"isLocked"`
}

type VariableDeleteInput struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

type VariableMutationOutcome struct {
	Errors   []*util2.Error        `json:"errors"`
	Variable *model.ConfigVariable `json:"variable"`
}
