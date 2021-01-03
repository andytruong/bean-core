package claim

import (
	"context"

	"bean/components/scalar"
)

// string -> *util.Payload
const ClaimsContextKey scalar.ContextKey = "bean.claims"

func PayloadToContext(ctx context.Context, payload *Payload) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, payload)
}

func ContextToPayload(ctx context.Context) *Payload {
	if claims, ok := ctx.Value(ClaimsContextKey).(*Payload); ok {
		return claims
	}

	return nil
}
