package namespace

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"bean/components/scalar"
	"bean/pkg/namespace/model"
	"bean/pkg/user"
	mUser "bean/pkg/user/model"
	"bean/pkg/util/connect"
)

func newMembershipResolver(bNamespace *NamespaceBean, bUser *user.UserBean) MembershipResolver {
	return MembershipResolver{
		bean: bNamespace,
		user: bUser,
	}
}

type MembershipResolver struct {
	bean *NamespaceBean
	user *user.UserBean
}

func (this MembershipResolver) Edges(ctx context.Context, obj *model.MembershipConnection) ([]*model.MembershipEdge, error) {
	var edges []*model.MembershipEdge

	for _, node := range obj.Nodes {
		edges = append(edges, &model.MembershipEdge{
			Cursor: model.MembershipNodeCursor(node),
			Node:   node,
		})
	}

	return edges, nil
}

func (this MembershipResolver) Namespace(ctx context.Context, obj *model.Membership) (*model.Namespace, error) {
	return this.bean.Load(ctx, obj.NamespaceID)
}

func (this MembershipResolver) User(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
	return this.user.Resolvers.Query.User(ctx, obj.UserID)
}

func (this MembershipResolver) UpdateLastLoginTime(db *gorm.DB, membership *model.Membership) error {
	membership.LoggedInAt = scalar.NilTime(time.Now())

	return db.Save(&membership).Error
}

func (this MembershipResolver) Roles(ctx context.Context, obj *model.Membership) ([]*model.Namespace, error) {
	return this.FindRoles(ctx, obj.UserID, obj.NamespaceID)
}

func (this MembershipResolver) FindRoles(ctx context.Context, userId string, namespaceId string) ([]*model.Namespace, error) {
	var roles []*model.Namespace

	err := this.bean.db.
		WithContext(ctx).
		Joins(
			fmt.Sprintf(
				"INNER JOIN %s ON %s.namespace_id = %s.id AND %s.user_id = ?",
				connect.TableNamespaceMemberships,
				connect.TableNamespaceMemberships,
				connect.TableNamespace,
				connect.TableNamespaceMemberships,
			),
			userId,
		).
		Where("kind = ?", model.NamespaceKindRole).
		Where("parent_id = ?", namespaceId).
		Find(&roles).
		Error

	if nil != err {
		return nil, err
	}

	return roles, nil
}
