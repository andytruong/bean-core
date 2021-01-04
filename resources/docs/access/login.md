Access › Login
====

1. Start session using credentials
====

Using credentials to get refresh token (`outcome.token`) to access system:

```graphql
mutation ($spaceId: ID!, $email: EmailAddress!) {
  accessMutation {
    session {
      create(
        input: {
            spaceId: $spaceId,
            email: $email,
            hashedPassword: "$hashedPassword",
            codeChallengeMethod: "S256",
            codeChallenge: "SHA256($codeVerifier)"
        }
      ) {
        errors { code fields message }
        token
      }
    }
  }
}
```

2. Access token
====

Go get access token (in JWT format) to access the system. It's fast to be expired, use session's `refreshToken`
and `$codeVerifier` to get it:

```graphql
query ($refreshToken: String!, $codeVerifier: String!) {
  accessQuery {
    session {
      load(token: $refreshToken) {
        id
        accessToken: jwt(codeVerifier: $codeVerifier)
      }
    }
  }
}
```

4. Use accessToken
====

```
curl http://path/to/endpoint -H "Authorization: Bearer ${accessToken}"
```

3. Start session using one-time login token
====

```graphql
mutation ($token: String!, $codeVerifier: String!) {
  accessMutation {
    session {
      generateOneTimeLoginToken(
        input: {
            token: $token
            codeChallengeMethod: "S256"
            codeChallenge: "SH256(…)"
        }
      ) {
        errors { code fields message }
        session { jwt(codeVerifier: $codeVerifier) }
        token
      }
    }
  }
}
```

4. Terminate the session
====

Related sessions will be also terminated (request with auth header).

```graphql
mutation {
    sessionArchive {
        errors  { code fields message }
        result
    }
}
```

References
====

- https://oauth.net/2/pkce/
- https://auth0.com/docs/flows/concepts/auth-code-pkce
