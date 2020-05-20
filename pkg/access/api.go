package access

import (
	"context"

	"bean/pkg/access/model/dto"
)

type AccessMutationResolver struct {
}

func (this *AccessMutationResolver) SessionCreate(ctx context.Context, input *dto.LoginInput) (*dto.LoginOutcome, error) {

	panic("not implemented")
}

func (this *AccessMutationResolver) SessionDelete(ctx context.Context, input *dto.LoginInput) (*dto.LogoutOutcome, error) {
	panic("not implemented")
}

type AccessQueryResolver struct {
}

func (this AccessQueryResolver) LoadSession(ctx context.Context, input *dto.ValidationInput) (*dto.ValidationOutcome, error) {
	panic("wip")
}
