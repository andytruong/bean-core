package infra

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	model1 "bean/pkg/app/model"
	dto1 "bean/pkg/app/model/dto"
	"bean/pkg/infra/gql"
	dto2 "bean/pkg/integration/mailer/model/dto"
	model2 "bean/pkg/integration/s3/model"
	dto3 "bean/pkg/integration/s3/model/dto"
	model3 "bean/pkg/space/model"
	dto4 "bean/pkg/space/model/dto"
	model4 "bean/pkg/user/model"
	dto5 "bean/pkg/user/model/dto"
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
func (r *Resolver) ApplicationMutation() gql.ApplicationMutationResolver {
	return &applicationMutationResolver{r}
}
func (r *Resolver) ApplicationQuery() gql.ApplicationQueryResolver {
	return &applicationQueryResolver{r}
}
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
func (r *Resolver) Mutation() gql.MutationResolver     { return &mutationResolver{r} }
func (r *Resolver) Query() gql.QueryResolver           { return &queryResolver{r} }
func (r *Resolver) S3Mutation() gql.S3MutationResolver { return &s3MutationResolver{r} }
func (r *Resolver) S3UploadMutation() gql.S3UploadMutationResolver {
	return &s3UploadMutationResolver{r}
}
func (r *Resolver) Session() gql.SessionResolver { return &sessionResolver{r} }
func (r *Resolver) Space() gql.SpaceResolver     { return &spaceResolver{r} }
func (r *Resolver) SpaceMembershipMutation() gql.SpaceMembershipMutationResolver {
	return &spaceMembershipMutationResolver{r}
}
func (r *Resolver) SpaceMembershipQuery() gql.SpaceMembershipQueryResolver {
	return &spaceMembershipQueryResolver{r}
}
func (r *Resolver) SpaceMutation() gql.SpaceMutationResolver { return &spaceMutationResolver{r} }
func (r *Resolver) SpaceQuery() gql.SpaceQueryResolver       { return &spaceQueryResolver{r} }
func (r *Resolver) User() gql.UserResolver                   { return &userResolver{r} }
func (r *Resolver) UserEmail() gql.UserEmailResolver         { return &userEmailResolver{r} }
func (r *Resolver) UserMutation() gql.UserMutationResolver   { return &userMutationResolver{r} }
func (r *Resolver) UserQuery() gql.UserQueryResolver         { return &userQueryResolver{r} }

// Resolvers
type accessMutationResolver struct{ *Resolver }
type accessQueryResolver struct{ *Resolver }
type accessSessionMutationResolver struct{ *Resolver }
type accessSessionQueryResolver struct{ *Resolver }
type applicationResolver struct{ *Resolver }
type applicationMutationResolver struct{ *Resolver }
type applicationQueryResolver struct{ *Resolver }
type mailerAccountMutationResolver struct{ *Resolver }
type mailerQueryAccountResolver struct{ *Resolver }
type mailerTemplateMutationResolver struct{ *Resolver }
type membershipResolver struct{ *Resolver }
type membershipConnectionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type s3MutationResolver struct{ *Resolver }
type s3UploadMutationResolver struct{ *Resolver }
type sessionResolver struct{ *Resolver }
type spaceResolver struct{ *Resolver }
type spaceMembershipMutationResolver struct{ *Resolver }
type spaceMembershipQueryResolver struct{ *Resolver }
type spaceMutationResolver struct{ *Resolver }
type spaceQueryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
type userEmailResolver struct{ *Resolver }
type userMutationResolver struct{ *Resolver }
type userQueryResolver struct{ *Resolver }

func (r *accessMutationResolver) Session(ctx context.Context, obj *dto.AccessMutation) (*dto.AccessSessionMutation, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["AccessMutation"].(map[string]interface{})
	callback := objectResolver["Session"].(func(ctx context.Context, obj *dto.AccessMutation) (*dto.AccessSessionMutation, error))

	return callback(ctx, obj)
}
func (r *accessQueryResolver) Session(ctx context.Context, obj *dto.AccessQuery) (*dto.AccessSessionQuery, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["AccessQuery"].(map[string]interface{})
	callback := objectResolver["Session"].(func(ctx context.Context, obj *dto.AccessQuery) (*dto.AccessSessionQuery, error))

	return callback(ctx, obj)
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
func (r *applicationResolver) Polices(ctx context.Context, obj *model1.Application) ([]*model2.Policy, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Application"].(map[string]interface{})
	callback := objectResolver["Polices"].(func(ctx context.Context, obj *model1.Application) ([]*model2.Policy, error))

	return callback(ctx, obj)
}
func (r *applicationResolver) Credentials(ctx context.Context, obj *model1.Application) (*model2.Credentials, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Application"].(map[string]interface{})
	callback := objectResolver["Credentials"].(func(ctx context.Context, obj *model1.Application) (*model2.Credentials, error))

	return callback(ctx, obj)
}
func (r *applicationMutationResolver) Create(ctx context.Context, obj *dto1.ApplicationMutation, input *dto1.ApplicationCreateInput) (*dto1.ApplicationOutcome, error) {
	panic("no implementation found in resolvers[ApplicationMutation][Create]")
}
func (r *applicationMutationResolver) Update(ctx context.Context, obj *dto1.ApplicationMutation, input *dto1.ApplicationUpdateInput) (*dto1.ApplicationOutcome, error) {
	panic("no implementation found in resolvers[ApplicationMutation][Update]")
}
func (r *applicationQueryResolver) Load(ctx context.Context, obj *dto1.ApplicationQuery, id string, version *string) (*model1.Application, error) {
	panic("no implementation found in resolvers[ApplicationQuery][Load]")
}
func (r *mailerAccountMutationResolver) Create(ctx context.Context, obj *dto2.MailerAccountMutation, input dto2.MailerAccountCreateInput) (*dto2.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto2.MailerAccountMutation, input dto2.MailerAccountCreateInput) (*dto2.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *mailerAccountMutationResolver) Update(ctx context.Context, obj *dto2.MailerAccountMutation, input dto2.MailerAccountUpdateInput) (*dto2.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto2.MailerAccountMutation, input dto2.MailerAccountUpdateInput) (*dto2.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *mailerAccountMutationResolver) Verify(ctx context.Context, obj *dto2.MailerAccountMutation, id string, version string) (*dto2.MailerAccountMutationOutcome, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MailerAccountMutation"].(map[string]interface{})
	callback := objectResolver["Verify"].(func(ctx context.Context, obj *dto2.MailerAccountMutation, id string, version string) (*dto2.MailerAccountMutationOutcome, error))

	return callback(ctx, obj, id, version)
}
func (r *mailerQueryAccountResolver) Get(ctx context.Context, obj *dto2.MailerQueryAccount, id string) (*dto2.MailerAccount, error) {
	panic("no implementation found in resolvers[MailerQueryAccount][Get]")
}
func (r *mailerQueryAccountResolver) GetMultiple(ctx context.Context, obj *dto2.MailerQueryAccount, first int, after *string) ([]*dto2.MailerAccount, error) {
	panic("no implementation found in resolvers[MailerQueryAccount][GetMultiple]")
}
func (r *mailerTemplateMutationResolver) Create(ctx context.Context, obj *dto2.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Create]")
}
func (r *mailerTemplateMutationResolver) Update(ctx context.Context, obj *dto2.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Update]")
}
func (r *mailerTemplateMutationResolver) Delete(ctx context.Context, obj *dto2.MailerTemplateMutation) (*bool, error) {
	panic("no implementation found in resolvers[MailerTemplateMutation][Delete]")
}
func (r *membershipResolver) Space(ctx context.Context, obj *model3.Membership) (*model3.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["Space"].(func(ctx context.Context, obj *model3.Membership) (*model3.Space, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) User(ctx context.Context, obj *model3.Membership) (*model4.User, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, obj *model3.Membership) (*model4.User, error))

	return callback(ctx, obj)
}
func (r *membershipResolver) Roles(ctx context.Context, obj *model3.Membership) ([]*model3.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Membership"].(map[string]interface{})
	callback := objectResolver["Roles"].(func(ctx context.Context, obj *model3.Membership) ([]*model3.Space, error))

	return callback(ctx, obj)
}
func (r *membershipConnectionResolver) Edges(ctx context.Context, obj *model3.MembershipConnection) ([]*model3.MembershipEdge, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["MembershipConnection"].(map[string]interface{})
	callback := objectResolver["Edges"].(func(ctx context.Context, obj *model3.MembershipConnection) ([]*model3.MembershipEdge, error))

	return callback(ctx, obj)
}
func (r *mutationResolver) AccessMutation(ctx context.Context) (*dto.AccessMutation, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["AccessMutation"].(func(ctx context.Context) (*dto.AccessMutation, error))

	return callback(ctx)
}
func (r *mutationResolver) ApplicationMutation(ctx context.Context) (*dto1.ApplicationMutation, error) {
	bundle, _ := r.container.bundles.App()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["ApplicationMutation"].(func(ctx context.Context) (*dto1.ApplicationMutation, error))

	return callback(ctx)
}
func (r *mutationResolver) MailerMutation(ctx context.Context) (*dto2.MailerMutation, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["MailerMutation"].(func(ctx context.Context) (*dto2.MailerMutation, error))

	return callback(ctx)
}
func (r *mutationResolver) S3Mutation(ctx context.Context) (*dto3.S3Mutation, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["S3Mutation"].(func(ctx context.Context) (*dto3.S3Mutation, error))

	return callback(ctx)
}
func (r *mutationResolver) SpaceMutation(ctx context.Context) (*dto4.SpaceMutation, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["SpaceMutation"].(func(ctx context.Context) (*dto4.SpaceMutation, error))

	return callback(ctx)
}
func (r *mutationResolver) UserMutation(ctx context.Context) (*dto5.UserMutation, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Mutation"].(map[string]interface{})
	callback := objectResolver["UserMutation"].(func(ctx context.Context) (*dto5.UserMutation, error))

	return callback(ctx)
}
func (r *queryResolver) AccessQuery(ctx context.Context) (*dto.AccessQuery, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["AccessQuery"].(func(ctx context.Context) (*dto.AccessQuery, error))

	return callback(ctx)
}
func (r *queryResolver) ApplicationQuery(ctx context.Context) (*dto1.ApplicationQuery, error) {
	bundle, _ := r.container.bundles.App()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["ApplicationQuery"].(func(ctx context.Context) (*dto1.ApplicationQuery, error))

	return callback(ctx)
}
func (r *queryResolver) MailerQuery(ctx context.Context) (*dto2.MailerQuery, error) {
	bundle, _ := r.container.bundles.Mailer()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["MailerQuery"].(func(ctx context.Context) (*dto2.MailerQuery, error))

	return callback(ctx)
}
func (r *queryResolver) SpaceQuery(ctx context.Context) (*dto4.SpaceQuery, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["SpaceQuery"].(func(ctx context.Context) (*dto4.SpaceQuery, error))

	return callback(ctx)
}
func (r *queryResolver) UserQuery(ctx context.Context) (*dto5.UserQuery, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Query"].(map[string]interface{})
	callback := objectResolver["UserQuery"].(func(ctx context.Context) (*dto5.UserQuery, error))

	return callback(ctx)
}
func (r *s3MutationResolver) Upload(ctx context.Context, obj *dto3.S3Mutation) (*dto3.S3UploadMutation, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["S3Mutation"].(map[string]interface{})
	callback := objectResolver["Upload"].(func(ctx context.Context, obj *dto3.S3Mutation) (*dto3.S3UploadMutation, error))

	return callback(ctx, obj)
}
func (r *s3UploadMutationResolver) Token(ctx context.Context, obj *dto3.S3UploadMutation, input dto3.S3UploadTokenInput) (map[string]interface{}, error) {
	bundle, _ := r.container.bundles.S3()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["S3UploadMutation"].(map[string]interface{})
	callback := objectResolver["Token"].(func(ctx context.Context, obj *dto3.S3UploadMutation, input dto3.S3UploadTokenInput) (map[string]interface{}, error))

	return callback(ctx, obj, input)
}
func (r *sessionResolver) User(ctx context.Context, obj *model.Session) (*model4.User, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["User"].(func(ctx context.Context, obj *model.Session) (*model4.User, error))

	return callback(ctx, obj)
}
func (r *sessionResolver) Space(ctx context.Context, obj *model.Session) (*model3.Space, error) {
	bundle, _ := r.container.bundles.Access()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Session"].(map[string]interface{})
	callback := objectResolver["Space"].(func(ctx context.Context, obj *model.Session) (*model3.Space, error))

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
func (r *spaceResolver) DomainNames(ctx context.Context, obj *model3.Space) (*model3.DomainNames, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["DomainNames"].(func(ctx context.Context, obj *model3.Space) (*model3.DomainNames, error))

	return callback(ctx, obj)
}
func (r *spaceResolver) Features(ctx context.Context, obj *model3.Space) (*model3.SpaceFeatures, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["Features"].(func(ctx context.Context, obj *model3.Space) (*model3.SpaceFeatures, error))

	return callback(ctx, obj)
}
func (r *spaceResolver) Parent(ctx context.Context, obj *model3.Space) (*model3.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["Space"].(map[string]interface{})
	callback := objectResolver["Parent"].(func(ctx context.Context, obj *model3.Space) (*model3.Space, error))

	return callback(ctx, obj)
}
func (r *spaceMembershipMutationResolver) Create(ctx context.Context, obj *dto4.SpaceMembershipMutation, input dto4.SpaceMembershipCreateInput) (*dto4.SpaceMembershipCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMembershipMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto4.SpaceMembershipMutation, input dto4.SpaceMembershipCreateInput) (*dto4.SpaceMembershipCreateOutcome, error))

	return callback(ctx, obj, input)
}
func (r *spaceMembershipMutationResolver) Update(ctx context.Context, obj *dto4.SpaceMembershipMutation, input dto4.SpaceMembershipUpdateInput) (*dto4.SpaceMembershipCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMembershipMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto4.SpaceMembershipMutation, input dto4.SpaceMembershipUpdateInput) (*dto4.SpaceMembershipCreateOutcome, error))

	return callback(ctx, obj, input)
}
func (r *spaceMembershipQueryResolver) Load(ctx context.Context, obj *dto4.SpaceMembershipQuery, id string, version *string) (*model3.Membership, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMembershipQuery"].(map[string]interface{})
	callback := objectResolver["Load"].(func(ctx context.Context, obj *dto4.SpaceMembershipQuery, id string, version *string) (*model3.Membership, error))

	return callback(ctx, obj, id, version)
}
func (r *spaceMembershipQueryResolver) Find(ctx context.Context, obj *dto4.SpaceMembershipQuery, first int, after *string, filters dto4.MembershipsFilter) (*model3.MembershipConnection, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMembershipQuery"].(map[string]interface{})
	callback := objectResolver["Find"].(func(ctx context.Context, obj *dto4.SpaceMembershipQuery, first int, after *string, filters dto4.MembershipsFilter) (*model3.MembershipConnection, error))

	return callback(ctx, obj, first, after, filters)
}
func (r *spaceMutationResolver) Create(ctx context.Context, obj *dto4.SpaceMutation, input dto4.SpaceCreateInput) (*dto4.SpaceCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto4.SpaceMutation, input dto4.SpaceCreateInput) (*dto4.SpaceCreateOutcome, error))

	return callback(ctx, obj, input)
}
func (r *spaceMutationResolver) Update(ctx context.Context, obj *dto4.SpaceMutation, input dto4.SpaceUpdateInput) (*dto4.SpaceCreateOutcome, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto4.SpaceMutation, input dto4.SpaceUpdateInput) (*dto4.SpaceCreateOutcome, error))

	return callback(ctx, obj, input)
}
func (r *spaceMutationResolver) Membership(ctx context.Context, obj *dto4.SpaceMutation) (*dto4.SpaceMembershipMutation, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceMutation"].(map[string]interface{})
	callback := objectResolver["Membership"].(func(ctx context.Context, obj *dto4.SpaceMutation) (*dto4.SpaceMembershipMutation, error))

	return callback(ctx, obj)
}
func (r *spaceQueryResolver) FindOne(ctx context.Context, obj *dto4.SpaceQuery, filters dto4.SpaceFilters) (*model3.Space, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceQuery"].(map[string]interface{})
	callback := objectResolver["FindOne"].(func(ctx context.Context, obj *dto4.SpaceQuery, filters dto4.SpaceFilters) (*model3.Space, error))

	return callback(ctx, obj, filters)
}
func (r *spaceQueryResolver) Membership(ctx context.Context, obj *dto4.SpaceQuery) (*dto4.SpaceMembershipQuery, error) {
	bundle, _ := r.container.bundles.Space()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["SpaceQuery"].(map[string]interface{})
	callback := objectResolver["Membership"].(func(ctx context.Context, obj *dto4.SpaceQuery) (*dto4.SpaceMembershipQuery, error))

	return callback(ctx, obj)
}
func (r *userResolver) Name(ctx context.Context, obj *model4.User) (*model4.UserName, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["User"].(map[string]interface{})
	callback := objectResolver["Name"].(func(ctx context.Context, obj *model4.User) (*model4.UserName, error))

	return callback(ctx, obj)
}
func (r *userResolver) Emails(ctx context.Context, obj *model4.User) (*model4.UserEmails, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["User"].(map[string]interface{})
	callback := objectResolver["Emails"].(func(ctx context.Context, obj *model4.User) (*model4.UserEmails, error))

	return callback(ctx, obj)
}
func (r *userEmailResolver) Verified(ctx context.Context, obj *model4.UserEmail) (bool, error) {
	panic("no implementation found in resolvers[UserEmail][Verified]")
}
func (r *userMutationResolver) Create(ctx context.Context, obj *dto5.UserMutation, input *dto5.UserCreateInput) (*dto5.UserMutationOutcome, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["UserMutation"].(map[string]interface{})
	callback := objectResolver["Create"].(func(ctx context.Context, obj *dto5.UserMutation, input *dto5.UserCreateInput) (*dto5.UserMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *userMutationResolver) Update(ctx context.Context, obj *dto5.UserMutation, input dto5.UserUpdateInput) (*dto5.UserMutationOutcome, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["UserMutation"].(map[string]interface{})
	callback := objectResolver["Update"].(func(ctx context.Context, obj *dto5.UserMutation, input dto5.UserUpdateInput) (*dto5.UserMutationOutcome, error))

	return callback(ctx, obj, input)
}
func (r *userQueryResolver) Load(ctx context.Context, obj *dto5.UserQuery, id string) (*model4.User, error) {
	bundle, _ := r.container.bundles.User()
	resolvers := bundle.GraphqlResolver()
	objectResolver := resolvers["UserQuery"].(map[string]interface{})
	callback := objectResolver["Load"].(func(ctx context.Context, obj *dto5.UserQuery, id string) (*model4.User, error))

	return callback(ctx, obj, id)
}
