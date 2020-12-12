package space

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"bean/components/scalar"
	"bean/pkg/space/model"
	"bean/pkg/user"
	mUser "bean/pkg/user/model"
	"bean/pkg/util/connect"
)

func newMembershipResolver(spaceBundle *SpaceBundle, userBundle *user.UserBundle) MembershipResolver {
	return MembershipResolver{
		bean: spaceBundle,
		user: userBundle,
	}
}

type MembershipResolver struct {
	bean *SpaceBundle
	user *user.UserBundle
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

func (this MembershipResolver) Space(ctx context.Context, obj *model.Membership) (*model.Space, error) {
	return this.bean.Load(ctx, obj.SpaceID)
}

func (this MembershipResolver) User(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
	return this.user.Resolvers.Query.User(ctx, obj.UserID)
}

func (this MembershipResolver) UpdateLastLoginTime(db *gorm.DB, membership *model.Membership) error {
	membership.LoggedInAt = scalar.NilTime(time.Now())

	return db.Save(&membership).Error
}

func (this MembershipResolver) Roles(ctx context.Context, obj *model.Membership) ([]*model.Space, error) {
	return this.FindRoles(ctx, obj.UserID, obj.SpaceID)
}

func (this MembershipResolver) FindRoles(ctx context.Context, userId string, spaceId string) ([]*model.Space, error) {
	var roles []*model.Space

	err := this.bean.db.
		WithContext(ctx).
		Joins(
			fmt.Sprintf(
				"INNER JOIN %s ON %s.space_id = %s.id AND %s.user_id = ?",
				connect.TableSpaceMemberships,
				connect.TableSpaceMemberships,
				connect.TableSpace,
				connect.TableSpaceMemberships,
			),
			userId,
		).
		Where("kind = ?", model.SpaceKindRole).
		Where("parent_id = ?", spaceId).
		Find(&roles).
		Error

	if nil != err {
		return nil, err
	}

	return roles, nil
}
