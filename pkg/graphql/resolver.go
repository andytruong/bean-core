package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	dto1 "bean/pkg/access/dto"
	"bean/pkg/user/dto"
	"context"
)

type Resolver struct{}

func (r *mutationResolver) Version(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto1.UserCreateOutcome, error) {
	panic("not implemented")
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	panic("not implemented")
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
