GraphQL API Standards
====

## Schema Design

Design your schema based on how data is used, not based on how it's stored.

- Field names should use camelCase. Many GraphQL clients are written in JavaScript, Java, Kotlin, or Swift, all of which recommend camelCase for variable names.
- Type names should use PascalCase. This matches how classes are defined in the languages mentioned above.
- Enum names should use PascalCase.
- Enum values should use ALL_CAPS, because they are similar to constants.

## Pagination

https://facebook.github.io/relay/graphql/connections.htm

## Errors

Mutations are likely to return errors and thus the response Payload must return an array of errors to inform the client 
about validation/processing errors

```
type MutationError {
    field: [String!]
    code: String!
    message: String!
}

type CreateLinkPayload {
  link: Link
  errors: [MutationError!]!
}
```

## Mutations

Mutations should be named as verbs, their inputs are the name with "Input" appended at the end.

For example a mutation to create a link in a bookmarks menu should take an input type of CreateLinkInput and return a 
response of created object.

Every mutation's response is to include the data that the mutation modified. This enables clients to obtain the latest 
persisted data without needing to send a followup query.

```
mutation CreateLinkMutation($input: CreateLinkInput!) {
  createLink(input: $input) {
    link {
      id
      createdAt
      url
      description
    }
  }
}
```

## More

- https://graphql-rules.com/
