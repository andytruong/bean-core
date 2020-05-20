package infra

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	model2 "bean/pkg/access/model"
	dto2 "bean/pkg/access/model/dto"
	"bean/pkg/infra/gql"
	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	model1 "bean/pkg/user/model"
	dto1 "bean/pkg/user/model/dto"
	"context"
)

type Resolver struct{}

func (r *mutationResolver) Ping(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *mutationResolver) NamespaceCreate(ctx context.Context, input dto.NamespaceCreateInput) (*dto.NamespaceCreateOutcome, error) {
	panic("not implemented")
}

func (r *mutationResolver) UserCreate(ctx context.Context, input *dto1.UserCreateInput) (*dto1.UserCreateOutcome, error) {
	panic("not implemented")
}

func (r *mutationResolver) SessionCreate(ctx context.Context, input *dto2.SessionCreateInput) (*dto2.SessionCreateOutcome, error) {
	panic("not implemented")
}

func (r *mutationResolver) SessionDelete(ctx context.Context, input *dto2.SessionCreateInput) (*dto2.LogoutOutcome, error) {
	panic("not implemented")
}

func (r *namespaceResolver) DomainNames(ctx context.Context, obj *model.Namespace) (*model.DomainNames, error) {
	panic("not implemented")
}

func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	panic("not implemented")
}

func (r *queryResolver) Namespace(ctx context.Context, id string) (*model.Namespace, error) {
	panic("not implemented")
}

func (r *queryResolver) User(ctx context.Context, id string) (*model1.User, error) {
	panic("not implemented")
}

func (r *queryResolver) LoadSession(ctx context.Context, input *dto2.ValidationInput) (*dto2.ValidationOutcome, error) {
	panic("not implemented")
}

func (r *sessionResolver) User(ctx context.Context, obj *model2.Session) (*model1.User, error) {
	panic("not implemented")
}

func (r *sessionResolver) Namespace(ctx context.Context, obj *model2.Session) (*model.Namespace, error) {
	panic("not implemented")
}

func (r *userResolver) Name(ctx context.Context, obj *model1.User) (*model1.UserName, error) {
	panic("not implemented")
}

func (r *userResolver) Emails(ctx context.Context, obj *model1.User) (*model1.UserEmails, error) {
	panic("not implemented")
}

func (r *userEmailResolver) Verified(ctx context.Context, obj *model1.UserEmail) (bool, error) {
	panic("not implemented")
}

// Mutation returns gql.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// Namespace returns gql.NamespaceResolver implementation.
func (r *Resolver) Namespace() gql.NamespaceResolver { return &namespaceResolver{r} }

// Query returns gql.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

// Session returns gql.SessionResolver implementation.
func (r *Resolver) Session() gql.SessionResolver { return &sessionResolver{r} }

// User returns gql.UserResolver implementation.
func (r *Resolver) User() gql.UserResolver { return &userResolver{r} }

// UserEmail returns gql.UserEmailResolver implementation.
func (r *Resolver) UserEmail() gql.UserEmailResolver { return &userEmailResolver{r} }

type mutationResolver struct{ *Resolver }
type namespaceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sessionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userEmailResolver struct{ *Resolver }
