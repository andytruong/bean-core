package infra

import (
	"context"

	"bean/pkg/access"
	"bean/pkg/namespace"
	"bean/pkg/user"
)

type queryResolver struct {
	user.UserQueryResolver
	*namespace.NamespaceQueryResolver
	*access.AccessQueryResolver
}

func (this queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
