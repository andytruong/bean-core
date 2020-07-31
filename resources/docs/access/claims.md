Claims
====

## Payload schema

```json
{
	"aud": "Namespace ID",
	"jti": "Session ID",
	"sub": "User ID",
	"exp": "ExpiresAt int64",
	"iat": "IssuedAt  int64",
	"iss": "Issuer",
	"nbf": "NotBefore int64",
	"kind": "credentials/onetime/authenticated/password-forget/password-reset"
	"roles": ["Role", "of", "user", "inside", "Namespace"]
}
```
