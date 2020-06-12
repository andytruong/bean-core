package namespace

import (
	"context"

	"bean/pkg/namespace/model/dto"
)

type API struct {
	Mutation MutationAPI
	Query    Query
}

type MutationAPI interface {
	CoreAPI
	MembershipAPI
}

type CoreAPI interface {
	NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error)
	NamespaceUpdate(ctx context.Context, input dto.NamespaceUpdateInput) (*bool, error)
}

type MembershipAPI interface {
	NamespaceMembershipCreate(ctx context.Context, input dto.NamespaceMembershipCreateInput) (*dto.NamespaceMembershipCreateOutcome, error)
	NamespaceMembershipUpdate(ctx context.Context, input dto.NamespaceMembershipUpdateInput) (*dto.NamespaceMembershipCreateOutcome, error)
}

type Query struct {
}
