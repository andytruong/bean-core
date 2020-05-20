package infra

import (
	"context"

	"bean/pkg/access"
	"bean/pkg/namespace"
	"bean/pkg/user"
)

type queryResolver struct {
	*user.UserModule
	*namespace.NamespaceModule
	*access.AccessModule
}

func (this queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
