package infra

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	model3 "bean/pkg/access/model"
	dto1 "bean/pkg/access/model/dto"
	"bean/pkg/infra/gql"
	"bean/pkg/integration/mailer/model/dto"
	"bean/pkg/integration/s3/model"
	dto4 "bean/pkg/integration/s3/model/dto"
	model1 "bean/pkg/space/model"
	dto2 "bean/pkg/space/model/dto"
	model2 "bean/pkg/user/model"
	dto3 "bean/pkg/user/model/dto"
	"context"
)

// Root resolver

type Resolver struct {
	container *Container
}

func (r *Resolver) Application() gql.ApplicationResolver { return &applicationResolver{r} }
func (r *Resolver) MailerAccountMutation() gql.MailerAccountMutationResolver {
	return &mailerAccountMutationResolver{r}
}
func (r *Resolver) MailerQueryAccount() gql.MailerQueryAccountResolver {
	return &mailerQueryAccountResolver{r}
}
func (r *Resolver) MailerTemplateMutation() gql.MailerTemplateMutationResolver {
	return &mailerTemplateMutationResolver{r}
}
func (r *Resolver) Membership() gql.MembershipResolver { return &membershipResolver{r} }
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
type mailerAccountMutationResolver struct{ *Resolver }
type mailerQueryAccountResolver struct{ *Resolver }
type mailerTemplateMutationResolver struct{ *Resolver }
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
func (r *mailerAccountMutationResolver) Create(ctx context.Context, obj *dto.MailerAccountMutation, input dto.MailerAccountCreateInput) (*dto.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto.MailerAccountMutation, input dto.MailerAccountCreateInput) (*dto.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *mailerAccountMutationResolver) Update(ctx context.Context, obj *dto.MailerAccountMutation, input dto.MailerAccountUpdateInput) (*dto.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto.MailerAccountMutation, input dto.MailerAccountUpdateInput) (*dto.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *mailerAccountMutationResolver) Verify(ctx context.Context, obj *dto.MailerAccountMutation, id string, version string) (*dto.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Verify"].(func(ctx context.Context, obj *dto.MailerAccountMutation, id string, version string) (*dto.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, id, version)
}
func (r *mailerQueryAccountResolver) Get(ctx context.Context, obj *dto.MailerQueryAccount, id string) (*dto.MailerAccount, error) {
	panic("no implementation found in resolvers[MailerQueryAccount][Get]")
}
func (r *mailerQueryAccountResolver) GetMultiple(ctx context.Context, obj *dto.MailerQueryAccount, first int, after *string) ([]*dto.MailerAccount, error) {
	panic("no implementation found in resolvers[MailerQueryAccount][GetMultiple]")
}
func (r *mailerTemplateMutationResolver) Create(ctx context.Context, obj *dto.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Create]")
}
func (r *mailerTemplateMutationResolver) Update(ctx context.Context, obj *dto.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Update]")
}
func (r *mailerTemplateMutationResolver) Delete(ctx context.Context, obj *dto.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Delete]")
}
func (r *membershipResolver) Space(ctx context.Context, obj *model1.Membership) (*model1.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["Space"].(func(ctx context.Context, obj *model1.Membership) (*model1.Space, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) User(ctx context.Context, obj *model1.Membership) (*model2.User, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, obj *model1.Membership) (*model2.User, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) Roles(ctx context.Context, obj *model1.Membership) ([]*model1.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["Roles"].(func(ctx context.Context, obj *model1.Membership) ([]*model1.Space, error))

	return callback(ctx, obj)
}
func (r *membershipConnectionResolver) Edges(ctx context.Context, obj *model1.MembershipConnection) ([]*model1.MembershipEdge, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MembershipConnection"].(map[string]interface{})
	callback := objectResolver["Edges"].(func(ctx context.Context, obj *model1.MembershipConnection) ([]*model1.MembershipEdge, error))

	return callback(ctx, obj)
}
func (r *mutationResolver) SessionCreate(ctx context.Context, input *dto1.SessionCreateInput) (*dto1.SessionCreateOutcome, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SessionCreate"].(func(ctx context.Context, input *dto1.SessionCreateInput) (*dto1.SessionCreateOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) SessionArchive(ctx context.Context) (*dto1.SessionArchiveOutcome, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SessionArchive"].(func(ctx context.Context) (*dto1.SessionArchiveOutcome, error))

	return callback(ctx)
}
func (r *mutationResolver) SpaceCreate(ctx context.Context, input dto2.SpaceCreateInput) (*dto2.SpaceCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SpaceCreate"].(func(ctx context.Context, input dto2.SpaceCreateInput) (*dto2.SpaceCreateOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) SpaceUpdate(ctx context.Context, input dto2.SpaceUpdateInput) (*dto2.SpaceCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SpaceUpdate"].(func(ctx context.Context, input dto2.SpaceUpdateInput) (*dto2.SpaceCreateOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) SpaceMembershipCreate(ctx context.Context, input dto2.SpaceMembershipCreateInput) (*dto2.SpaceMembershipCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SpaceMembershipCreate"].(func(ctx context.Context, input dto2.SpaceMembershipCreateInput) (*dto2.SpaceMembershipCreateOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) SpaceMembershipUpdate(ctx context.Context, input dto2.SpaceMembershipUpdateInput) (*dto2.SpaceMembershipCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SpaceMembershipUpdate"].(func(ctx context.Context, input dto2.SpaceMembershipUpdateInput) (*dto2.SpaceMembershipCreateOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) UserCreate(ctx context.Context, input *dto3.UserCreateInput) (*dto3.UserMutationOutcome, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["UserCreate"].(func(ctx context.Context, input *dto3.UserCreateInput) (*dto3.UserMutationOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) UserUpdate(ctx context.Context, input dto3.UserUpdateInput) (*dto3.UserMutationOutcome, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["UserUpdate"].(func(ctx context.Context, input dto3.UserUpdateInput) (*dto3.UserMutationOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) S3ApplicationCreate(ctx context.Context, input *dto4.S3ApplicationCreateInput) (*dto4.S3ApplicationMutationOutcome, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3ApplicationCreate"].(func(ctx context.Context, input *dto4.S3ApplicationCreateInput) (*dto4.S3ApplicationMutationOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) S3ApplicationUpdate(ctx context.Context, input *dto4.S3ApplicationUpdateInput) (*dto4.S3ApplicationMutationOutcome, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3ApplicationUpdate"].(func(ctx context.Context, input *dto4.S3ApplicationUpdateInput) (*dto4.S3ApplicationMutationOutcome, error))

	return callback(ctx, input)
}
func (r *mutationResolver) S3UploadToken(ctx context.Context, input dto4.S3UploadTokenInput) (map[string]interface{}, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3UploadToken"].(func(ctx context.Context, input dto4.S3UploadTokenInput) (map[string]interface{}, error))

	return callback(ctx, input)
}
func (r *mutationResolver) MailerMutation(ctx context.Context) (*dto.MailerMutation, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["MailerMutation"].(func(ctx context.Context) (*dto.MailerMutation, error))

	return callback(ctx)
}
func (r *queryResolver) Session(ctx context.Context, token string) (*model3.Session, error) {
	panic("no implementation found in resolvers[Query][Session]")
}
func (r *queryResolver) Space(ctx context.Context, filters dto2.SpaceFilters) (*model1.Space, error) {
	panic("no implementation found in resolvers[Query][Space]")
}
func (r *queryResolver) Membership(ctx context.Context, id string, version *string) (*model1.Membership, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["Membership"].(func(ctx context.Context, id string, version *string) (*model1.Membership, error))

	return callback(ctx, id, version)
}
func (r *queryResolver) Memberships(ctx context.Context, first int, after *string, filters dto2.MembershipsFilter) (*model1.MembershipConnection, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["Memberships"].(func(ctx context.Context, first int, after *string, filters dto2.MembershipsFilter) (*model1.MembershipConnection, error))

	return callback(ctx, first, after, filters)
}
func (r *queryResolver) User(ctx context.Context, id string) (*model2.User, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, id string) (*model2.User, error))

	return callback(ctx, id)
}
func (r *queryResolver) MailerQuery(ctx context.Context) (*dto.MailerQuery, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["MailerQuery"].(func(ctx context.Context) (*dto.MailerQuery, error))

	return callback(ctx)
}
func (r *sessionResolver) User(ctx context.Context, obj *model3.Session) (*model2.User, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, obj *model3.Session) (*model2.User, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Space(ctx context.Context, obj *model3.Session) (*model1.Space, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Space"].(func(ctx context.Context, obj *model3.Session) (*model1.Space, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Scopes(ctx context.Context, obj *model3.Session) ([]*model3.AccessScope, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Scopes"].(func(ctx context.Context, obj *model3.Session) ([]*model3.AccessScope, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Context(ctx context.Context, obj *model3.Session) (*model3.SessionContext, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Context"].(func(ctx context.Context, obj *model3.Session) (*model3.SessionContext, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Jwt(ctx context.Context, obj *model3.Session, codeVerifier string) (string, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Jwt"].(func(ctx context.Context, obj *model3.Session, codeVerifier string) (string, error))

	return callback(ctx, obj, codeVerifier)
}
func (r *spaceResolver) DomainNames(ctx context.Context, obj *model1.Space) (*model1.DomainNames, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["DomainNames"].(func(ctx context.Context, obj *model1.Space) (*model1.DomainNames, error))

	return callback(ctx, obj)
}
func (r *spaceResolver) Features(ctx context.Context, obj *model1.Space) (*model1.SpaceFeatures, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["Features"].(func(ctx context.Context, obj *model1.Space) (*model1.SpaceFeatures, error))

	return callback(ctx, obj)
}
func (r *spaceResolver) Parent(ctx context.Context, obj *model1.Space) (*model1.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["Parent"].(func(ctx context.Context, obj *model1.Space) (*model1.Space, error))

	return callback(ctx, obj)
}
func (r *userResolver) Name(ctx context.Context, obj *model2.User) (*model2.UserName, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["User"].(map[string]interface{})
	callback := objectResolver["Name"].(func(ctx context.Context, obj *model2.User) (*model2.UserName, error))

	return callback(ctx, obj)
}
func (r *userResolver) Emails(ctx context.Context, obj *model2.User) (*model2.UserEmails, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["User"].(map[string]interface{})
	callback := objectResolver["Emails"].(func(ctx context.Context, obj *model2.User) (*model2.UserEmails, error))

	return callback(ctx, obj)
}
func (r *userEmailResolver) Verified(ctx context.Context, obj *model2.UserEmail) (bool, error) {
	panic("no implementation found in resolvers[UserEmail][Verified]")
}
