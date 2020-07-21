package namespace

import (
	"context"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

func newResolver(bean *NamespaceBean) *Resolvers {
	return &Resolvers{
		Object:   &NamespaceObjectResolver{bean: bean},
		Query:    &NamespaceQueryResolver{bean: bean},
		Mutation: &NamespaceMutationResolver{bean: bean},
	}
}

type (
	Resolvers struct {
		Object   *NamespaceObjectResolver
		Query    *NamespaceQueryResolver
		Mutation *NamespaceMutationResolver
	}

	NamespaceObjectResolver struct {
		bean *NamespaceBean
	}

	NamespaceQueryResolver struct {
		bean *NamespaceBean
	}
	NamespaceMutationResolver struct {
		bean *NamespaceBean
	}
)

func (this NamespaceObjectResolver) Parent(ctx context.Context, obj *model.Namespace) (*model.Namespace, error) {
	if nil == obj.ParentID {
		return nil, nil
	}

	return this.bean.Load(ctx, *obj.ParentID)
}

func (this NamespaceObjectResolver) DomainNames(ctx context.Context, namespace *model.Namespace) (*model.DomainNames, error) {
	return this.bean.CoreDomainName.Find(namespace)
}

func (this NamespaceObjectResolver) Features(ctx context.Context, namespace *model.Namespace) (*model.NamespaceFeatures, error) {
	return this.bean.CoreConfig.List(ctx, namespace)
}

func (this NamespaceQueryResolver) Memberships(ctx context.Context, first int, after *string, filters dto.MembershipsFilter) (*model.MembershipConnection, error) {
	return this.bean.CoreMember.Find(first, after, filters)
}

func (this NamespaceQueryResolver) Membership(ctx context.Context, id string, version *string) (*model.Membership, error) {
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
