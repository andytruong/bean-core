package access

import (
	"context"

	"bean/components/claim"
	"bean/components/util"
	"bean/pkg/access/model"
	"bean/pkg/access/model/dto"
	space_model "bean/pkg/space/model"
	user_model "bean/pkg/user/model"
)

func (bundle *Bundle) newResolves() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"AccessQuery": func(ctx context.Context) (*dto.AccessQuery, error) {
				return &dto.AccessQuery{}, nil
			},
		},
		"AccessQuery": map[string]interface{}{
			"Session": func(ctx context.Context, _ *dto.AccessQuery) (*dto.AccessSessionQuery, error) {
				return &dto.AccessSessionQuery{}, nil
			},
		},
		"AccessSessionQuery": map[string]interface{}{
			"Load": func(ctx context.Context, _ *dto.AccessSessionQuery, id string) (*model.Session, error) {
				return bundle.sessionService.load(ctx, id)
			},
		},
		"Mutation": map[string]interface{}{
			"AccessMutation": func(ctx context.Context) (*dto.AccessMutation, error) {
				return &dto.AccessMutation{}, nil
			},
		},
		"AccessMutation": map[string]interface{}{
			"Session": func(ctx context.Context, _ *dto.AccessMutation) (*dto.AccessSessionMutation, error) {
				return &dto.AccessSessionMutation{}, nil
			},
		},
		"AccessSessionMutation": map[string]interface{}{
			"Create": func(ctx context.Context, _ *dto.AccessSessionMutation, in *dto.SessionCreateInput) (
				*dto.SessionOutcome, error,
			) {
				return bundle.sessionService.newSessionWithCredentials(ctx, in)
			},
			"GenerateOneTimeLoginToken": func(ctx context.Context, _ *dto.AccessSessionMutation, in *dto.SessionCreateOTLTSessionInput) (
				*dto.SessionOutcome, error,
			) {
				return bundle.sessionService.newOTLTSession(ctx, in)
			},
			"ExchangeOneTimeLoginToken": func(ctx context.Context, _ *dto.AccessSessionMutation, in *dto.SessionExchangeOTLTInput) (
				*dto.SessionOutcome, error,
			) {
				return bundle.sessionService.newSessionWithOTLT(ctx, in)
			},
			"Archive": func(ctx context.Context, _ *dto.AccessSessionMutation) (*dto.SessionArchiveOutcome, error) {
				claims := claim.ContextToPayload(ctx)
				if nil == claims {
					return nil, util.ErrorAuthRequired
				}

				sess, _ := bundle.sessionService.load(ctx, claims.SessionId())
				if sess != nil {
					out, err := bundle.sessionService.Delete(ctx, sess)

					return out, err
				}

				return &dto.SessionArchiveOutcome{
					Errors: util.NewErrors(util.ErrorCodeInput, []string{"token"}, "session not found"),
					Result: false,
				}, nil
			},
		},
		"Session": map[string]interface{}{
			"User": func(ctx context.Context, obj *model.Session) (*user_model.User, error) {
				return bundle.userBundle.UserService.Load(ctx, obj.UserId)
			},
			"Context": func(ctx context.Context, obj *model.Session) (*model.SessionContext, error) {
				panic("implement me")
			},
			"Scopes": func(ctx context.Context, obj *model.Session) ([]*model.AccessScope, error) {
				return obj.Scopes, nil
			},
			"Space": func(ctx context.Context, obj *model.Session) (*space_model.Space, error) {
				return bundle.spaceBundle.Service.Load(ctx, obj.SpaceId)
			},
			"Jwt": func(ctx context.Context, session *model.Session, codeVerifier string) (string, error) {
				return bundle.JwtService.getSignedString(ctx, session, codeVerifier)
			},
		},
	}
}
