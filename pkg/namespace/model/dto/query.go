package dto

import (
	"bean/pkg/util"
)

type NamespaceFilters struct {
	ID     *string   `json:"id"`
	Domain *util.Uri `json:"domain"`
}
