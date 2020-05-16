package namespace

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

type (
	NamespaceMutationResolver struct {
		db *gorm.DB
		id *util.Identifier
	}

	NamespaceQueryResolver struct {
		db *gorm.DB
	}

	NamespaceModelResolver struct {
		db *gorm.DB
	}
)

func (this NamespaceMutationResolver) NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	panic("implement me: NamespaceMutationResolver.NamespaceCreate()")
}

func (this NamespaceModelResolver) DomainNames(ctx context.Context, obj *model.Namespace) (*model.DomainNames, error) {
	panic("implement me")
}
