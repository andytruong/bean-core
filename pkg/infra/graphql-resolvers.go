package infra

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	model3 "bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/infra/gql"
	dto4 "bean/pkg/integration/mailer/model/dto"
	"bean/pkg/integration/s3/model"
	dto3 "bean/pkg/integration/s3/model/dto"
	model1 "bean/pkg/space/model"
	dto1 "bean/pkg/space/model/dto"
	model2 "bean/pkg/user/model"
	dto2 "bean/pkg/user/model/dto"
	"context"
)

// Root resolver

type Resolver struct {
	container Container
}

func (r *Resolver) Application() gql.ApplicationResolver { return &applicationResolver{r} }
func (r *Resolver) Membership() gql.MembershipResolver   { return &membershipResolver{r} }
func (r *Resolver) MembershipConnection() gql.MembershipConnectionResolver {
	return &membershipConnectionResolver{r}
}
func (r *Resolver) Mutation() gql.MutationResolver   { return &mutationResolver{r} }
func (r *Resolver) Query() gql.QueryResolver         { return &queryResolver{r} }
func (r *Resolver) Session() gql.SessionResolver     { return &sessionResolver{r} }
func (r *Resolver) Space() gql.SpaceResolver         { return &spaceResolver{r} }
func (r *Resolver) User() gql.UserResolver           { return &userResolver{r} }
func (r *Resolver) UserEmail() gql.UserEmailResolver { return &userEmailResolver{r} }

// Resolvers
type applicationResolver struct{ *Resolver }
type membershipResolver struct{ *Resolver }
type membershipConnectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sessionResolver struct{ *Resolver }
type spaceResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userEmailResolver struct{ *Resolver }

func (r *applicationResolver) Polices(ctx context.Context, obj *model.Application) ([]*model.Policy, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Application"].(map[string]interface{})
	callback := objectResolver["Polices"].(func(ctx context.Context, obj *model.Application) ([]*model.Policy, error))

	return callback(ctx, obj)
}
func (r *applicationResolver) Credentials(ctx context.Context, obj *model.Application) (*model.Credentials, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Application"].(map[string]interface{})
	callback := objectResolver["Credentials"].(func(ctx context.Context, obj *model.Application) (*model.Credentials, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) Space(ctx context.Context, obj *model1.Membership) (*model1.Space, error) {
	panic("no implementation")
}
func (r *membershipResolver) User(ctx context.Context, obj *model1.Membership) (*model2.User, error) {
	panic("no implementation")
}
func (r *membershipResolver) Roles(ctx context.Context, obj *model1.Membership) ([]*model1.Space, error) {
	panic("no implementation")
}
func (r *membershipConnectionResolver) Edges(ctx context.Context, obj *model1.MembershipConnection) ([]*model1.MembershipEdge, error) {
	panic("no implementation")
}
func (r *mutationResolver) SessionCreate(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SessionCreate"].(func(ctx context.Context, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) SessionArchive(ctx context.Context) (*dto.SessionArchiveOutcome, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SessionArchive"].(func(ctx context.Context) (*dto.SessionArchiveOutcome, error))

	return callback(ctx)
}
func (r *mutationResolver) SpaceMembershipCreate(ctx context.Context, input dto1.SpaceMembershipCreateInput) (*dto1.SpaceMembershipCreateOutcome, error) {
	panic("no implementation")
}
func (r *mutationResolver) SpaceMembershipUpdate(ctx context.Context, input dto1.SpaceMembershipUpdateInput) (*dto1.SpaceMembershipCreateOutcome, error) {
	panic("no implementation")
}
func (r *mutationResolver) SpaceCreate(ctx context.Context, input dto1.SpaceCreateInput) (*dto1.SpaceCreateOutcome, error) {
	panic("no implementation")
}
func (r *mutationResolver) SpaceUpdate(ctx context.Context, input dto1.SpaceUpdateInput) (*dto1.SpaceCreateOutcome, error) {
	panic("no implementation")
}
func (r *mutationResolver) UserCreate(ctx context.Context, input *dto2.UserCreateInput) (*dto2.UserMutationOutcome, error) {
	panic("no implementation")
}
func (r *mutationResolver) UserUpdate(ctx context.Context, input dto2.UserUpdateInput) (*dto2.UserMutationOutcome, error) {
	panic("no implementation")
}
func (r *mutationResolver) S3ApplicationCreate(ctx context.Context, input *dto3.S3ApplicationCreateInput) (*dto3.S3ApplicationMutationOutcome, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3ApplicationCreate"].(func(ctx context.Context, input *dto3.S3ApplicationCreateInput) (*dto3.S3ApplicationMutationOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) S3ApplicationUpdate(ctx context.Context, input *dto3.S3ApplicationUpdateInput) (*dto3.S3ApplicationMutationOutcome, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3ApplicationUpdate"].(func(ctx context.Context, input *dto3.S3ApplicationUpdateInput) (*dto3.S3ApplicationMutationOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) S3UploadToken(ctx context.Context, input dto3.S3UploadTokenInput) (map[string]interface{}, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3UploadToken"].(func(ctx context.Context, input dto3.S3UploadTokenInput) (map[string]interface{}, error))

	return callback(ctx, input)
}
func (r *mutationResolver) MailerMutation(ctx context.Context) (*dto4.MailerMutation, error) {
	panic("no implementation")
}
func (r *queryResolver) Session(ctx context.Context, token string) (*model3.Session, error) {
	panic("no implementation")
}
func (r *queryResolver) Membership(ctx context.Context, id string, version *string) (*model1.Membership, error) {
	panic("no implementation")
}
func (r *queryResolver) Memberships(ctx context.Context, first int, after *string, filters dto1.MembershipsFilter) (*model1.MembershipConnection, error) {
	panic("no implementation")
}
func (r *queryResolver) Space(ctx context.Context, filters dto1.SpaceFilters) (*model1.Space, error) {
	panic("no implementation")
}
func (r *queryResolver) User(ctx context.Context, id string) (*model2.User, error) {
	panic("no implementation")
}
func (r *queryResolver) MailerQuery(ctx context.Context) (*dto4.MailerQuery, error) {
	panic("no implementation")
}
func (r *sessionResolver) User(ctx context.Context, obj *model3.Session) (*model2.User, error) {
	panic("no implementation")
}
func (r *sessionResolver) Space(ctx context.Context, obj *model3.Session) (*model1.Space, error) {
	panic("no implementation")
}
func (r *sessionResolver) Scopes(ctx context.Context, obj *model3.Session) ([]*model3.AccessScope, error) {
	panic("no implementation")
}
func (r *sessionResolver) Context(ctx context.Context, obj *model3.Session) (*model3.SessionContext, error) {
	panic("no implementation")
}
func (r *sessionResolver) Jwt(ctx context.Context, obj *model3.Session, codeVerifier string) (string, error) {
	panic("no implementation")
}
func (r *spaceResolver) DomainNames(ctx context.Context, obj *model1.Space) (*model1.DomainNames, error) {
	panic("no implementation")
}
func (r *spaceResolver) Features(ctx context.Context, obj *model1.Space) (*model1.SpaceFeatures, error) {
	panic("no implementation")
}
func (r *spaceResolver) Parent(ctx context.Context, obj *model1.Space) (*model1.Space, error) {
	panic("no implementation")
}
func (r *userResolver) Name(ctx context.Context, obj *model2.User) (*model2.UserName, error) {
	panic("no implementation")
}
func (r *userResolver) Emails(ctx context.Context, obj *model2.User) (*model2.UserEmails, error) {
	panic("no implementation")
}
func (r *userEmailResolver) Verified(ctx context.Context, obj *model2.UserEmail) (bool, error) {
	panic("no implementation")
}
