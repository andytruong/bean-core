Claims
====

## Payload schema

```json
{
	"aud": "Namespace ID",
	"exp": "ExpiresAt int64",
	"jti": "Session ID",
	"iat": "IssuedAt  int64",
	"iss": "Issuer",
	"nbf": "NotBefore int64",
	"sub": "User ID",
	"kind": "credentials/onetime/authenticated/password-forget/password-reset"
	"roles": ["Role", "of", "user", "inside", "Namespace"]
}
```
