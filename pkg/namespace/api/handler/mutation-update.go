package handler

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model/dto"
)

type NamespaceModuleHandler struct {
	DB *gorm.DB
}

func (this NamespaceModuleHandler) NamespaceUpdate(ctx context.Context, input dto.NamespaceUpdateInput) (*bool, error) {
	// load the namespace
	// check version for conflict
	// update
	// change version
	panic("implement me")
}
