package gql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	access_dto "bean/pkg/access/dto"
	"bean/pkg/infra"
	"bean/pkg/user/dto"
)

func NewResolver(container *infra.Container) (ResolverRoot, error) {
	resolver := &Resolver{
		container: container,
	}

	return resolver, nil
}

type Resolver struct {
	container *infra.Container
}

func (r *mutationResolver) Version(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*access_dto.UserCreateOutcome, error) {
	panic("not implemented")
}

func (r *mutationResolver) SessionCreate(ctx context.Context, input *access_dto.LoginInput) (*access_dto.LoginOutcome, error) {
	panic("not implemented")
}

func (r *mutationResolver) SessionDelete(ctx context.Context, input *access_dto.LoginInput) (*access_dto.LogoutPayload, error) {
	panic("not implemented")
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) LoadSession(ctx context.Context, input *access_dto.ValidationInput) (*access_dto.ValidationOutcome, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
