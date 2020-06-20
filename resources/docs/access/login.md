Access › Login
====

## Login using credentials

Using credentials to get access-token (outcome.session.jwt) and refresh token (outcome.token) to access system:

```graphql
mutation ($namespaceId: ID!, $email: EmailAddress!) {
    sessionCreate(
        input: {
            credentials: {
                namespaceId: $namespaceId,
                email: $email
                hashedPassword: String!
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

## Login using one-time login token

```graphql
mutation (token: String!) {
    sessionCreate(
        input: {
            oneTimeLogin: {
                token: $token
            }
        }
    ) {
        errors  { code fields message }
        session { jwt }
        token
    }
}
```
