package namespace

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/pkg/namespace/fixtures"
	"bean/pkg/util"
)

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	input := fixtures.NamespaceCreateInputFixture()

	t.Run("happy case", func(t *testing.T) {
		db := util.MockDatabase()
		module, err := NewNamespaceModule(db, util.MockLogger(), util.MockIdentifier())
		ass.NoError(err)
		util.MockInstall(module, db)

		now := time.Now()
		outcome, err := module.Mutation.NamespaceCreate(context.Background(), input)
		ass.NoError(err)
		ass.Nil(outcome.Errors)
		ass.Equal(input.Object.Title, outcome.Namespace.Title)
		ass.Equal(input.Object.IsActive, outcome.Namespace.IsActive)
		ass.True(outcome.Namespace.CreatedAt.UnixNano() >= now.UnixNano())
		ass.True(outcome.Namespace.UpdatedAt.UnixNano() >= now.UnixNano())
	})
}

func Test_Create_Error(t *testing.T) {
	t.Run("domain duplication", func(t *testing.T) {
		ass := assert.New(t)
		input := fixtures.NamespaceCreateInputFixture()

		db := util.MockDatabase()
		module, err := NewNamespaceModule(db, util.MockLogger(), util.MockIdentifier())
		ass.NoError(err)
		util.MockInstall(module, db)

		// create first
		{
			outcome, err := module.Mutation.NamespaceCreate(context.Background(), input)
			ass.NoError(err)
			ass.Nil(outcome.Errors)
		}

		// create again with same input
		{
			outcome, err := module.Mutation.NamespaceCreate(context.Background(), input)
			ass.NoError(err)

			fmt.Println("err", err, outcome)
			// ass.NoError(err)
			// ass.Nil(outcome.Errors)
		}
	})
}
