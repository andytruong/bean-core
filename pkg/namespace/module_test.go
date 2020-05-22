package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/pkg/namespace/api/fixtures"
	"bean/pkg/util"
)

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	input := fixtures.NamespaceCreateInputFixture()
	db := util.MockDatabase()
	module := NewNamespaceModule(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(module, db)

	t.Run("happy case", func(t *testing.T) {
		now := time.Now()
		outcome, err := module.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		ass.Nil(outcome.Errors)
		ass.Equal(input.Object.Title, outcome.Namespace.Title)
		ass.Equal(input.Object.IsActive, outcome.Namespace.IsActive)
		ass.True(outcome.Namespace.CreatedAt.UnixNano() >= now.UnixNano())
		ass.True(outcome.Namespace.UpdatedAt.UnixNano() >= now.UnixNano())
	})

	t.Run("domain duplication", func(t *testing.T) {
		// create again with same input
		outcome, err := module.NamespaceCreate(context.Background(), input)

		ass.Nil(outcome)
		ass.NotNil(err)
		ass.Contains(err.Error(), "UNIQUE constraint failed: namespace_domains.value")
	})
}

func Test_Query(t *testing.T) {
	ass := assert.New(t)
	db := util.MockDatabase()
	module := NewNamespaceModule(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(module, db)

	var id string
	input := fixtures.NamespaceCreateInputFixture()

	{
		// setup data for query
		outcome, err := module.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		id = outcome.Namespace.ID
	}

	{
		obj, err := module.Namespace(context.Background(), id)
		ass.NoError(err)
		ass.Equal(obj.ID, id)
		ass.Equal(obj.Title, input.Object.Title)
		ass.Equal(obj.IsActive, input.Object.IsActive)
	}
}
