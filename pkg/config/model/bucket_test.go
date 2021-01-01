package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config_Bucket(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)

	t.Run("schema validation", func(t *testing.T) {
		t.Run("example", func(t *testing.T) {
			bucket := ConfigBucket{
				Schema: `{
    "title": "Person",
    "type": "object",
    "properties": {
        "firstName": { "type": "string" },
        "lastName":  { "type": "string" },
        "age":       { "type": "integer", "minimum": 0 },
        "friends": {
          "type" : "array",
          "items" : { "title" : "REFERENCE", "$ref" : "#" }
        }
    },
    "required": ["firstName", "lastName"]
  }`,
			}

			t.Run("invalid value", func(t *testing.T) {
				reasons, err := bucket.Validate(ctx, `{}`)

				ass.Len(reasons, 2)
				ass.Contains(reasons[0], `"firstName" value is required`)
				ass.Contains(reasons[1], `"lastName" value is required`)
				ass.Nil(err)
			})

			t.Run("valid value", func(t *testing.T) {
				reasons, err := bucket.Validate(ctx, `{"firstName": "John", "lastName": "Doe"}`)

				ass.Nil(reasons)
				ass.Nil(err)
			})
		})

		t.Run("simple", func(t *testing.T) {
			bucket := ConfigBucket{Schema: `{"type": "number"}`}

			reasons, err := bucket.Validate(ctx, `"OK"`)
			ass.Nil(err)
			ass.Contains(reasons[0], `"OK" type should be number, got string`)
		})
	})
}
