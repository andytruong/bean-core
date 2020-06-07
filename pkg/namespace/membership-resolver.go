package namespace

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"

	"bean/pkg/namespace/model"
	"bean/pkg/user"
	mUser "bean/pkg/user/model"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

func newMembershipResolver(namespaceModule *NamespaceModule, userModule *user.UserModule) MembershipResolver {
	return MembershipResolver{
		namespace: namespaceModule,
		user:      userModule,
	}
}

type MembershipResolver struct {
	namespace *NamespaceModule
	user      *user.UserModule
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
	return this.namespace.Namespace(ctx, obj.NamespaceID)
}

func (this MembershipResolver) User(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
	return this.user.User(ctx, obj.UserID)
}

func (this MembershipResolver) UpdateLastLoginTime(db *gorm.DB, membership *model.Membership) error {
	membership.LoggedInAt = util.NilTime(time.Now())

	return db.
		Table(connect.TableNamespaceMemberships).
		Save(&membership).
		Error
}
