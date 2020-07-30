Access › Login
====

1. Start session using credentials
====

Using credentials to get access-token (outcome.session.jwt) and refresh token (outcome.token) to access system:

```graphql
mutation ($namespaceId: ID!, $email: EmailAddress!) {
    sessionCreate(
        input: {
            credentials: {
                namespaceId: $namespaceId,
                email: $email
                hashedPassword: String!
                codeChallengeMethod: "S256",
                codeChallenge: "SHA256(…)"
            }
        }
    ) {
        errors  { code fields message }

        # JWT can be used to access system, expire in 5m
        # -------
        session { jwt }

        # Token can be used to refresh the token.
        # -------
        token
    }
}
```

2. Refresh token
====

When JWT is expired, we can refresh it

```graphql
query ($token: String!, $codeVerifier: String!) {
    session(token: $token) {
        jwt($codeVerifier)
    }
}
```

3. Start session using one-time login token
====

```graphql
mutation (token: String!) {
    sessionCreate(
        input: {
            oneTimeLogin: { token: $token }
            codeChallengeMethod: "S256",
            codeChallenge: "SHA256(…)"
        }
    ) {
        errors  { code fields message }
        session { jwt }
        token
    }
}
```

4. Terminate the session
====

```graphql
# TODO
```

References
====

- https://oauth.net/2/pkce/
- https://auth0.com/docs/flows/concepts/auth-code-pkce
