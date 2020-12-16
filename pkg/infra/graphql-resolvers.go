package infra

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	"bean/pkg/infra/gql"
	dto1 "bean/pkg/integration/mailer/model/dto"
	model1 "bean/pkg/integration/s3/model"
	dto4 "bean/pkg/integration/s3/model/dto"
	model2 "bean/pkg/space/model"
	dto2 "bean/pkg/space/model/dto"
	model3 "bean/pkg/user/model"
	dto3 "bean/pkg/user/model/dto"
	"context"
)

// Root resolver

type Resolver struct {
	container *Container
}

func (r *Resolver) AccessMutation() gql.AccessMutationResolver { return &accessMutationResolver{r} }
func (r *Resolver) AccessQuery() gql.AccessQueryResolver       { return &accessQueryResolver{r} }
func (r *Resolver) AccessSessionMutation() gql.AccessSessionMutationResolver {
	return &accessSessionMutationResolver{r}
}
func (r *Resolver) AccessSessionQuery() gql.AccessSessionQueryResolver {
	return &accessSessionQueryResolver{r}
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
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }
func (r *Resolver) Query() gql.QueryResolver       { return &queryResolver{r} }
func (r *Resolver) S3ApplicationMutation() gql.S3ApplicationMutationResolver {
	return &s3ApplicationMutationResolver{r}
}
func (r *Resolver) S3Mutation() gql.S3MutationResolver { return &s3MutationResolver{r} }
func (r *Resolver) S3UploadMutation() gql.S3UploadMutationResolver {
	return &s3UploadMutationResolver{r}
}
func (r *Resolver) Session() gql.SessionResolver     { return &sessionResolver{r} }
func (r *Resolver) Space() gql.SpaceResolver         { return &spaceResolver{r} }
func (r *Resolver) User() gql.UserResolver           { return &userResolver{r} }
func (r *Resolver) UserEmail() gql.UserEmailResolver { return &userEmailResolver{r} }

// Resolvers
type accessMutationResolver struct{ *Resolver }
type accessQueryResolver struct{ *Resolver }
type accessSessionMutationResolver struct{ *Resolver }
type accessSessionQueryResolver struct{ *Resolver }
type applicationResolver struct{ *Resolver }
type mailerAccountMutationResolver struct{ *Resolver }
type mailerQueryAccountResolver struct{ *Resolver }
type mailerTemplateMutationResolver struct{ *Resolver }
type membershipResolver struct{ *Resolver }
type membershipConnectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type s3ApplicationMutationResolver struct{ *Resolver }
type s3MutationResolver struct{ *Resolver }
type s3UploadMutationResolver struct{ *Resolver }
type sessionResolver struct{ *Resolver }
type spaceResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userEmailResolver struct{ *Resolver }

func (r *accessMutationResolver) Session(ctx context.Context, obj *dto.AccessMutation) (*dto.AccessSessionMutation, error) {
	panic("no implementation found in resolvers[AccessMutation][Session]")
}
func (r *accessQueryResolver) Session(ctx context.Context, obj *dto.AccessQuery) (*dto.AccessSessionQuery, error) {
	panic("no implementation found in resolvers[AccessQuery][Session]")
}
func (r *accessSessionMutationResolver) Create(ctx context.Context, obj *dto.AccessSessionMutation, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["AccessSessionMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto.AccessSessionMutation, input *dto.SessionCreateInput) (*dto.SessionCreateOutcome, error))

	return callback(ctx, obj, input)
}
func (r *accessSessionMutationResolver) Archive(ctx context.Context, obj *dto.AccessSessionMutation) (*dto.SessionArchiveOutcome, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["AccessSessionMutation"].(map[string]interface{})
	callback := objectResolver["Archive"].(func(ctx context.Context, obj *dto.AccessSessionMutation) (*dto.SessionArchiveOutcome, error))

	return callback(ctx, obj)
}
func (r *accessSessionQueryResolver) Load(ctx context.Context, obj *dto.AccessSessionQuery, token string) (*model.Session, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["AccessSessionQuery"].(map[string]interface{})
	callback := objectResolver["Load"].(func(ctx context.Context, obj *dto.AccessSessionQuery, token string) (*model.Session, error))

	return callback(ctx, obj, token)
}
func (r *applicationResolver) Polices(ctx context.Context, obj *model1.Application) ([]*model1.Policy, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Application"].(map[string]interface{})
	callback := objectResolver["Polices"].(func(ctx context.Context, obj *model1.Application) ([]*model1.Policy, error))

	return callback(ctx, obj)
}
func (r *applicationResolver) Credentials(ctx context.Context, obj *model1.Application) (*model1.Credentials, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Application"].(map[string]interface{})
	callback := objectResolver["Credentials"].(func(ctx context.Context, obj *model1.Application) (*model1.Credentials, error))

	return callback(ctx, obj)
}
func (r *mailerAccountMutationResolver) Create(ctx context.Context, obj *dto1.MailerAccountMutation, input dto1.MailerAccountCreateInput) (*dto1.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto1.MailerAccountMutation, input dto1.MailerAccountCreateInput) (*dto1.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *mailerAccountMutationResolver) Update(ctx context.Context, obj *dto1.MailerAccountMutation, input dto1.MailerAccountUpdateInput) (*dto1.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto1.MailerAccountMutation, input dto1.MailerAccountUpdateInput) (*dto1.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *mailerAccountMutationResolver) Verify(ctx context.Context, obj *dto1.MailerAccountMutation, id string, version string) (*dto1.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Verify"].(func(ctx context.Context, obj *dto1.MailerAccountMutation, id string, version string) (*dto1.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, id, version)
}
func (r *mailerQueryAccountResolver) Get(ctx context.Context, obj *dto1.MailerQueryAccount, id string) (*dto1.MailerAccount, error) {
	panic("no implementation found in resolvers[MailerQueryAccount][Get]")
}
func (r *mailerQueryAccountResolver) GetMultiple(ctx context.Context, obj *dto1.MailerQueryAccount, first int, after *string) ([]*dto1.MailerAccount, error) {
	panic("no implementation found in resolvers[MailerQueryAccount][GetMultiple]")
}
func (r *mailerTemplateMutationResolver) Create(ctx context.Context, obj *dto1.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Create]")
}
func (r *mailerTemplateMutationResolver) Update(ctx context.Context, obj *dto1.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Update]")
}
func (r *mailerTemplateMutationResolver) Delete(ctx context.Context, obj *dto1.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Delete]")
}
func (r *membershipResolver) Space(ctx context.Context, obj *model2.Membership) (*model2.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["Space"].(func(ctx context.Context, obj *model2.Membership) (*model2.Space, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) User(ctx context.Context, obj *model2.Membership) (*model3.User, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, obj *model2.Membership) (*model3.User, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) Roles(ctx context.Context, obj *model2.Membership) ([]*model2.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["Roles"].(func(ctx context.Context, obj *model2.Membership) ([]*model2.Space, error))

	return callback(ctx, obj)
}
func (r *membershipConnectionResolver) Edges(ctx context.Context, obj *model2.MembershipConnection) ([]*model2.MembershipEdge, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MembershipConnection"].(map[string]interface{})
	callback := objectResolver["Edges"].(func(ctx context.Context, obj *model2.MembershipConnection) ([]*model2.MembershipEdge, error))

	return callback(ctx, obj)
}
func (r *mutationResolver) AccessMutation(ctx context.Context) (*dto.AccessMutation, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["AccessMutation"].(func(ctx context.Context) (*dto.AccessMutation, error))

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
func (r *mutationResolver) S3Mutation(ctx context.Context) (*dto4.S3Mutation, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3Mutation"].(func(ctx context.Context) (*dto4.S3Mutation, error))

	return callback(ctx)
}
func (r *mutationResolver) MailerMutation(ctx context.Context) (*dto1.MailerMutation, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["MailerMutation"].(func(ctx context.Context) (*dto1.MailerMutation, error))

	return callback(ctx)
}
func (r *queryResolver) AccessQuery(ctx context.Context) (*dto.AccessQuery, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["AccessQuery"].(func(ctx context.Context) (*dto.AccessQuery, error))

	return callback(ctx)
}
func (r *queryResolver) Space(ctx context.Context, filters dto2.SpaceFilters) (*model2.Space, error) {
	panic("no implementation found in resolvers[Query][Space]")
}
func (r *queryResolver) Membership(ctx context.Context, id string, version *string) (*model2.Membership, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["Membership"].(func(ctx context.Context, id string, version *string) (*model2.Membership, error))

	return callback(ctx, id, version)
}
func (r *queryResolver) Memberships(ctx context.Context, first int, after *string, filters dto2.MembershipsFilter) (*model2.MembershipConnection, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["Memberships"].(func(ctx context.Context, first int, after *string, filters dto2.MembershipsFilter) (*model2.MembershipConnection, error))

	return callback(ctx, first, after, filters)
}
func (r *queryResolver) User(ctx context.Context, id string) (*model3.User, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, id string) (*model3.User, error))

	return callback(ctx, id)
}
func (r *queryResolver) MailerQuery(ctx context.Context) (*dto1.MailerQuery, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["MailerQuery"].(func(ctx context.Context) (*dto1.MailerQuery, error))

	return callback(ctx)
}
func (r *s3ApplicationMutationResolver) Create(ctx context.Context, obj *dto4.S3ApplicationMutation, input *dto4.S3ApplicationCreateInput) (*dto4.S3ApplicationMutationOutcome, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["S3ApplicationMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto4.S3ApplicationMutation, input *dto4.S3ApplicationCreateInput) (*dto4.S3ApplicationMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *s3ApplicationMutationResolver) Update(ctx context.Context, obj *dto4.S3ApplicationMutation, input *dto4.S3ApplicationUpdateInput) (*dto4.S3ApplicationMutationOutcome, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["S3ApplicationMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto4.S3ApplicationMutation, input *dto4.S3ApplicationUpdateInput) (*dto4.S3ApplicationMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *s3MutationResolver) Application(ctx context.Context, obj *dto4.S3Mutation) (*dto4.S3ApplicationMutation, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["S3Mutation"].(map[string]interface{})
	callback := objectResolver["Application"].(func(ctx context.Context, obj *dto4.S3Mutation) (*dto4.S3ApplicationMutation, error))

	return callback(ctx, obj)
}
func (r *s3MutationResolver) Uploadd(ctx context.Context, obj *dto4.S3Mutation) (*dto4.S3UploadMutation, error) {
	panic("no implementation found in resolvers[S3Mutation][Uploadd]")
}
func (r *s3UploadMutationResolver) Token(ctx context.Context, obj *dto4.S3UploadMutation, input dto4.S3UploadTokenInput) (map[string]interface{}, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["S3UploadMutation"].(map[string]interface{})
	callback := objectResolver["Token"].(func(ctx context.Context, obj *dto4.S3UploadMutation, input dto4.S3UploadTokenInput) (map[string]interface{}, error))

	return callback(ctx, obj, input)
}
func (r *sessionResolver) User(ctx context.Context, obj *model.Session) (*model3.User, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, obj *model.Session) (*model3.User, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Space(ctx context.Context, obj *model.Session) (*model2.Space, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Space"].(func(ctx context.Context, obj *model.Session) (*model2.Space, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Scopes(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Scopes"].(func(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Context(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Context"].(func(ctx context.Context, obj *model.Session) (*model.SessionContext, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Jwt(ctx context.Context, obj *model.Session, codeVerifier string) (string, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Jwt"].(func(ctx context.Context, obj *model.Session, codeVerifier string) (string, error))

	return callback(ctx, obj, codeVerifier)
}
func (r *spaceResolver) DomainNames(ctx context.Context, obj *model2.Space) (*model2.DomainNames, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["DomainNames"].(func(ctx context.Context, obj *model2.Space) (*model2.DomainNames, error))

	return callback(ctx, obj)
}
func (r *spaceResolver) Features(ctx context.Context, obj *model2.Space) (*model2.SpaceFeatures, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["Features"].(func(ctx context.Context, obj *model2.Space) (*model2.SpaceFeatures, error))

	return callback(ctx, obj)
}
func (r *spaceResolver) Parent(ctx context.Context, obj *model2.Space) (*model2.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["Parent"].(func(ctx context.Context, obj *model2.Space) (*model2.Space, error))

	return callback(ctx, obj)
}
func (r *userResolver) Name(ctx context.Context, obj *model3.User) (*model3.UserName, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["User"].(map[string]interface{})
	callback := objectResolver["Name"].(func(ctx context.Context, obj *model3.User) (*model3.UserName, error))

	return callback(ctx, obj)
}
func (r *userResolver) Emails(ctx context.Context, obj *model3.User) (*model3.UserEmails, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["User"].(map[string]interface{})
	callback := objectResolver["Emails"].(func(ctx context.Context, obj *model3.User) (*model3.UserEmails, error))

	return callback(ctx, obj)
}
func (r *userEmailResolver) Verified(ctx context.Context, obj *model3.UserEmail) (bool, error) {
	panic("no implementation found in resolvers[UserEmail][Verified]")
}
