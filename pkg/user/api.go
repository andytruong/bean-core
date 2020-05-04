package user

import (
	"context"

	"bean/pkg/user/dto"
)

type (
	UserMutationResolver struct {
	}
)

func (this *UserMutationResolver) UserCreate(ctx context.Context, input *dto.UserCreateInput) (*dto.UserCreateOutcome, error) {
	panic("not implemented")
}
