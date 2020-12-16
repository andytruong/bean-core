package space

import (
	"context"

	"github.com/pkg/errors"

	"bean/components/scalar"
	util2 "bean/components/util"
	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	mUser "bean/pkg/user/model"
)

func (this *SpaceBundle) newResolvers() map[string]interface{} {
	return map[string]interface{}{
		"Query": map[string]interface{}{
			"SpaceQuery": func(ctx context.Context) (*dto.SpaceQuery, error) {
				return &dto.SpaceQuery{}, nil
			},
		},
		"SpaceQuery": map[string]interface{}{
			"FindOne": func(ctx context.Context, filters dto.SpaceFilters) (*model.Space, error) {
				return this.Service.FindOne(ctx, filters)
			},
			"Membership": func(ctx context.Context) (*dto.SpaceMembershipQuery, error) {
				return &dto.SpaceMembershipQuery{}, nil
			},
		},
		"SpaceMembershipQuery": map[string]interface{}{
			"Find": func(ctx context.Context, first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
				return this.MemberService.Find(first, after, filters)
			},
			"Load": func(ctx context.Context, id string, version *string) (*model.Membership, error) {
				return this.MemberService.load(ctx, id, version)
			},
		},
		"Mutation": map[string]interface{}{
			"SpaceMutation": func(ctx context.Context) (*dto.SpaceMutation, error) {
				return &dto.SpaceMutation{}, nil
			},
		},
		"SpaceMutation": map[string]interface{}{
			"Create": func(ctx context.Context, input dto.SpaceCreateInput) (*dto.SpaceCreateOutcome, error) {
				txn := this.db.WithContext(ctx).Begin()
				out, err := this.Service.Create(txn, input)

				if nil != err {
					txn.Rollback()
					return nil, err
				} else {
					return out, txn.Commit().Error
				}
			},
			"Update": func(ctx context.Context, in dto.SpaceUpdateInput) (*dto.SpaceCreateOutcome, error) {
				space, err := this.Service.Load(ctx, in.SpaceID)
				if nil != err {
					return nil, err
				}

				txn := this.db.WithContext(ctx).Begin()
				out, err := this.Service.Update(txn, *space, in)

				if nil != err {
					txn.Rollback()
					return nil, err
				} else {
					return out, txn.Commit().Error
				}
			},
			"Membership": func(ctx context.Context) (*dto.SpaceMembershipMutation, error) {
				return &dto.SpaceMembershipMutation{}, nil
			},
		},
		"SpaceMembershipMutation": map[string]interface{}{
			"Create": func(ctx context.Context, in dto.SpaceMembershipCreateInput) (*dto.SpaceMembershipCreateOutcome, error) {
				space, err := this.Service.Load(ctx, in.SpaceID)
				if nil != err {
					return nil, err
				}

				_, err = this.userBundle.Service.Load(this.db.WithContext(ctx), in.UserID)
				if nil != err {
					return nil, err
				}

				features, err := this.configService.List(ctx, space)
				if nil != err {
					return nil, err
				}

				if !features.Register {
					return nil, errors.Wrap(util2.ErrorConfig, "register is off")
				}

				tx := this.db.WithContext(ctx).Begin()
				outcome, err := this.MemberService.Create(tx, in)

				if nil != err {
					tx.Rollback()
					return nil, err
				} else {
					return outcome, tx.Commit().Error
				}
			},
			"Update": func(ctx context.Context, in dto.SpaceMembershipUpdateInput) (*dto.SpaceMembershipCreateOutcome, error) {
				membership, err := this.MemberService.load(ctx, in.Id, scalar.NilString(in.Version))

				if nil != err {
					return nil, err
				}

				tx := this.db.WithContext(ctx).Begin()
				outcome, err := this.MemberService.Update(tx, in, membership)

				if nil != err {
					tx.Rollback()
					return nil, err
				} else {
					return outcome, tx.Commit().Error
				}
			},
		},
		"Space": map[string]interface{}{
			"Parent": func(ctx context.Context, obj *model.Space) (*model.Space, error) {
				if nil == obj.ParentID {
					return nil, nil
				}

				return this.Service.Load(ctx, *obj.ParentID)
			},
			"DomainNames": func(ctx context.Context, space *model.Space) (*model.DomainNames, error) {
				return this.domainNameService.Find(space)
			},
			"Features": func(ctx context.Context, space *model.Space) (*model.SpaceFeatures, error) {
				return this.configService.List(ctx, space)
			},
		},
		"Membership": map[string]interface{}{
			"Space": func(ctx context.Context, obj *model.Membership) (*model.Space, error) {
				return this.Service.Load(ctx, obj.SpaceID)
			},
			"User": func(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
				return this.userBundle.Service.Load(this.db.WithContext(ctx), obj.UserID)
			},
			"Roles": func(ctx context.Context, obj *model.Membership) ([]*model.Space, error) {
				return this.MemberService.FindRoles(ctx, obj.UserID, obj.SpaceID)
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