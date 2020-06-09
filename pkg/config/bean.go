package config

import (
	"context"

	"bean/pkg/config/model/dto"
)

type ConfigBean struct {
}

func (this ConfigBean) BucketCreate(ctx context.Context, input dto.BucketCreateInput) (*dto.BucketMutationOutcome, error) {
	panic("todo")
}

func (this ConfigBean) BucketUpdate(ctx context.Context, input dto.BucketUpdateInput) (*dto.BucketMutationOutcome, error) {
	panic("todo")
}

// API to create bucket
// Bean install
// API to update bucket
