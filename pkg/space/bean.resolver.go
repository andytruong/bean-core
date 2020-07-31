package space

import (
	"context"

	"bean/pkg/space/model"
	"bean/pkg/space/model/dto"
	"bean/pkg/util"
)

func newResolver(bean *SpaceBean) *Resolvers {
	return &Resolvers{
		Object:   &SpaceObjectResolver{bean: bean},
		Query:    &SpaceQueryResolver{bean: bean},
		Mutation: &SpaceMutationResolver{bean: bean},
	}
}

type (
	Resolvers struct {
		Object   *SpaceObjectResolver
		Query    *SpaceQueryResolver
		Mutation *SpaceMutationResolver
	}

	SpaceObjectResolver struct {
		bean *SpaceBean
	}

	SpaceQueryResolver struct {
		bean *SpaceBean
	}
	SpaceMutationResolver struct {
		bean *SpaceBean
	}
)

func (this SpaceObjectResolver) Parent(ctx context.Context, obj *model.Space) (*model.Space, error) {
	if nil == obj.ParentID {
		return nil, nil
	}

	return this.bean.Load(ctx, *obj.ParentID)
}

func (this SpaceObjectResolver) DomainNames(ctx context.Context, space *model.Space) (*model.DomainNames, error) {
	return this.bean.CoreDomainName.Find(space)
}

func (this SpaceObjectResolver) Features(ctx context.Context, space *model.Space) (*model.SpaceFeatures, error) {
	return this.bean.CoreConfig.List(ctx, space)
}

func (this SpaceQueryResolver) Memberships(ctx context.Context, first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
	return this.bean.CoreMember.Find(first, after, filters)
}

func (this SpaceQueryResolver) Membership(ctx context.Context, id string, version *string) (*model.Membership, error) {
	obj := &model.Membership{}

	err := this.bean.db.WithContext(ctx).First(&obj, "id = ?", id).Error
	if nil != err {
		return nil, err
	} else if nil != version {
		if obj.Version != *version {
			return nil, util.ErrorVersionConflict
		}
	}

	return obj, nil
}
