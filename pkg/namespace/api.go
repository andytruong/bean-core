package namespace

import (
	"context"

	"github.com/jinzhu/gorm"

	"bean/pkg/infra/gql"
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
)

func (n NamespaceMutationResolver) NamespaceCreate(ctx context.Context, input gql.NamespaceCreateInput) (*gql.NamespaceCreateOutcome, error) {
	panic("implement me: NamespaceMutationResolver.NamespaceCreate()")
}
