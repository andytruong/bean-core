package claim

import (
	"context"

	"bean/components/scalar"
)

// string -> *util.Payload
const ClaimsContextKey scalar.ContextKey = "bean.claims"

func ClaimPayloadToContext(ctx context.Context, payload *Payload) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, payload)
}
