package namespace

import (
	"context"

	"bean/pkg/namespace/model"
	"bean/pkg/user"
	mUser "bean/pkg/user/model"
)

func newMembershipResolver(namespaceModule *NamespaceModule, userModule *user.UserModule) MembershipResolver {
	return MembershipResolver{
		namespaceModule: namespaceModule,
		userModule:      userModule,
	}
}

type MembershipResolver struct {
	namespaceModule *NamespaceModule
	userModule      *user.UserModule
}

func (this MembershipResolver) Namespace(ctx context.Context, obj *model.Membership) (*model.Namespace, error) {
	return this.namespaceModule.Namespace(ctx, obj.NamespaceID)
}

func (this MembershipResolver) User(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
	return this.userModule.User(ctx, obj.UserID)
}
