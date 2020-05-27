package api

import "context"

type Resolver struct {
}

func (this *Resolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
