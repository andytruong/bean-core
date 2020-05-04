package access

import (
	"context"

	"bean/pkg/access/dto"
)

type AccessMutationResolver struct {
}

func (r *AccessMutationResolver) SessionCreate(ctx context.Context, input *dto.LoginInput) (*dto.LoginOutcome, error) {
	panic("not implemented")
}

func (r *AccessMutationResolver) SessionDelete(ctx context.Context, input *dto.LoginInput) (*dto.LogoutOutcome, error) {
	panic("not implemented")
}

type QueryResolver struct {
}

func (this QueryResolver) LoadSession(ctx context.Context, input *dto.ValidationInput) (*dto.ValidationOutcome, error) {
	panic("wip")
}
