package util

import "context"

type CtxKey string

const (
	// string -> *util.Claims
	CxtKeyClaims CtxKey = "bean.claims"
)

func (this CtxKey) Actor(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(this).(*Claims); ok {
		return claims
	}

	return nil
}
