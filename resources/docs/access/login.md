Access › Login
====

1. Start session using credentials
====

Using credentials to get refresh token (`outcome.token`) to access system:

```graphql
mutation ($namespaceId: ID!, $email: EmailAddress!) {
    sessionCreate(
        input: {
            credentials: {
                namespaceId: $namespaceId,
                email: $email
                hashedPassword: String!
                codeChallengeMethod: "S256",
                codeChallenge: "SHA256($codeVerifier)"
            }
        }
    ) {
        errors  { code fields message }

        # Token can be used to refresh the token.
        # -------
        token
    }
}
```

2. Access token
====

Go get access token (in JWT format) to access the system. It's fast to be expired, use session's `refreshToken`
and `$codeVerifier` to get it:

```graphql
query ($refreshToken: String!, $codeVerifier: String!) {
    session(token: $refreshToken) {
        accessToken: jwt($codeVerifier)
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

> TODO: To be implemented.

Related sessions will be also terminated.

```graphql
mutation {
    sessionDelete(token: $accessToken) {
        errors { code message }
    }
}
```

References
====

- https://oauth.net/2/pkce/
- https://auth0.com/docs/flows/concepts/auth-code-pkce
