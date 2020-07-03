package claim

// See resources/docs/access/claims.md for use cases
// maxLength: 32
type Kind string

const (
	KindCredentials    Kind = "credentials"
	KindOTLT           Kind = "onetime"
	KindAuthenticated  Kind = "authenticated"
	KindPasswordForgot Kind = "password-forget"
	KindPasswordReset  Kind = "password-reset"
)
