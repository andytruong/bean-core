GraphQL API
====

## 1. Define GraphQL schema

```graphql
# file: pkg/my_bundle/api/schema.graphql

extend Query {
    hello: String!
}
```

## 2. Update glgen.yml

```yaml
schema:
    # …
    - bean/pkg/my_bundle/api/*.graphql

autobind:
    # …
    - bean/pkg/my_bundle/model
    - bean/pkg/my_bundle/dto
```

## 3. Implement bundle interface

```go
type Bundle interface {
	// …
	GraphqlResolver() map[string]interface{}
}
```

Example

```go
func (this MyBundle) GraphqlResolver() map[string]interface{} {
    return map[string]interface{} {
        "Query": map[string]interface{} {
            "Hello": func(ctx context.Context) (string, error) {
                return "Hello world!", nil
            }
        }
    }
}
```

## 4. Regenerate

```
make gql
```
