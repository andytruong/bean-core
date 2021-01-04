package space

import (
	"context"

	"github.com/pkg/errors"

	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	mUser "bean/pkg/user/model"
)

func (bundle *Bundle) newResolvers() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"SpaceQuery": func(ctx context.Context) (*dto.SpaceQuery, error) {
				return &dto.SpaceQuery{}, nil
			},
		},
		"SpaceQuery": map[string]interface{}{
			"FindOne": func(ctx context.Context, _ *dto.SpaceQuery, filters dto.SpaceFilters) (*model.Space, error) {
				return bundle.Service.FindOne(ctx, filters)
			},
			"Membership": func(ctx context.Context, _ *dto.SpaceQuery) (*dto.SpaceMembershipQuery, error) {
				return &dto.SpaceMembershipQuery{}, nil
			},
		},
		"SpaceMembershipQuery": map[string]interface{}{
			"Find": func(ctx context.Context, _ *dto.SpaceMembershipQuery, first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
				return bundle.MemberService.Find(ctx, first, after, filters)
			},
			"Load": func(ctx context.Context, _ *dto.SpaceMembershipQuery, id string, version *string) (*model.Membership, error) {
				return bundle.MemberService.load(ctx, id, version)
			},
		},
		"Mutation": map[string]interface{}{
			"SpaceMutation": func(ctx context.Context) (*dto.SpaceMutation, error) {
				return &dto.SpaceMutation{}, nil
			},
		},
		"SpaceMutation": map[string]interface{}{
			"Create": func(ctx context.Context, _ *dto.SpaceMutation, input dto.SpaceCreateInput) (*dto.SpaceOutcome, error) {
				tx := connect.ContextToDB(ctx)
				out, err := bundle.Service.Create(ctx, input)

				if nil != err {
					tx.Rollback()
					return nil, err
				} else {
					return out, tx.Commit().Error
				}
			},
			"Update": func(ctx context.Context, _ *dto.SpaceMutation, in dto.SpaceUpdateInput) (*dto.SpaceOutcome, error) {
				space, err := bundle.Service.Load(ctx, in.SpaceID)
				if nil != err {
					return nil, err
				}

				txn := connect.ContextToDB(ctx)
				out, err := bundle.Service.Update(ctx, *space, in)

				if nil != err {
					txn.Rollback()
					return nil, err
				} else {
					return out, txn.Commit().Error
				}
			},
			"Membership": func(ctx context.Context, _ *dto.SpaceMutation) (*dto.SpaceMembershipMutation, error) {
				return &dto.SpaceMembershipMutation{}, nil
			},
		},
		"SpaceMembershipMutation": map[string]interface{}{
			"Create": func(ctx context.Context, _ *dto.SpaceMembershipMutation, in dto.SpaceMembershipCreateInput) (*dto.SpaceMembershipCreateOutcome, error) {
				space, err := bundle.Service.Load(ctx, in.SpaceID)
				if nil != err {
					return nil, err
				}

				_, err = bundle.userBundle.Service.Load(ctx, in.UserID)
				if nil != err {
					return nil, err
				}

				features, err := bundle.configService.List(ctx, space)
				if nil != err {
					return nil, err
				}

				if !features.Register {
					return nil, errors.Wrap(util.ErrorConfig, "register is off")
				}

				tx := connect.ContextToDB(ctx).Begin()

				outcome, err := bundle.MemberService.Create(connect.DBToContext(ctx, tx), in)

				if nil != err {
					tx.Rollback()
					return nil, err
				} else {
					return outcome, tx.Commit().Error
				}
			},
			"Update": func(ctx context.Context, _ *dto.SpaceMembershipMutation, in dto.SpaceMembershipUpdateInput) (*dto.SpaceMembershipCreateOutcome, error) {
				membership, err := bundle.MemberService.load(ctx, in.Id, scalar.NilString(in.Version))

				if nil != err {
					return nil, err
				}

				tx := connect.ContextToDB(ctx).Begin()
				out, err := bundle.MemberService.Update(connect.DBToContext(ctx, tx), in, membership)

				if nil != err {
					tx.Rollback()
					return nil, err
				} else {
					return out, tx.Commit().Error
				}
			},
		},
		"Space": map[string]interface{}{
			"Parent": func(ctx context.Context, obj *model.Space) (*model.Space, error) {
				if nil == obj.ParentID {
					return nil, nil
				}

				return bundle.Service.Load(ctx, *obj.ParentID)
			},
			"DomainNames": func(ctx context.Context, space *model.Space) (*model.DomainNames, error) {
				return bundle.domainNameService.Find(ctx, space)
			},
			"Features": func(ctx context.Context, space *model.Space) (*model.SpaceFeatures, error) {
				return bundle.configService.List(ctx, space)
			},
		},
		"Membership": map[string]interface{}{
			"Space": func(ctx context.Context, obj *model.Membership) (*model.Space, error) {
				return bundle.Service.Load(ctx, obj.SpaceID)
			},
			"User": func(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
				return bundle.userBundle.Service.Load(ctx, obj.UserID)
			},
			"Roles": func(ctx context.Context, obj *model.Membership) ([]*model.Space, error) {
				return bundle.MemberService.FindRoles(ctx, obj.UserID, obj.SpaceID)
			},
		},
		"MembershipConnection": map[string]interface{}{
			"Edges": func(ctx context.Context, obj *model.MembershipConnection) ([]*model.MembershipEdge, error) {
				var edges []*model.MembershipEdge

				for _, node := range obj.Nodes {
					edges = append(edges, &model.MembershipEdge{
						Cursor: model.MembershipNodeCursor(node),
						Node:   node,
					})
				}

				return edges, nil
			},
		},
	}
}
