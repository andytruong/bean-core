package namespace

import (
	"context"

	"bean/pkg/namespace/model"
	mUser "bean/pkg/user/model"
)

func newMembershipResolver() MembershipResolver {
	return MembershipResolver{}
}

type MembershipResolver struct{}

func (this MembershipResolver) Namespace(ctx context.Context, obj *model.Membership) (*model.Namespace, error) {
	panic("implement me")
}

func (this MembershipResolver) User(ctx context.Context, obj *model.Membership) (*mUser.User, error) {
	panic("implement me")
}
