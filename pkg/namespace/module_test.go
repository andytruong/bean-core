package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bean/pkg/namespace/api/fixtures"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util"
)

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	input := fixtures.NamespaceCreateInputFixture()
	db := util.MockDatabase().LogMode(false)
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

func Test_Update(t *testing.T) {
	ass := assert.New(t)
	input := fixtures.NamespaceCreateInputFixture()
	db := util.MockDatabase()
	module := NewNamespaceModule(db, util.MockLogger(), util.MockIdentifier())
	util.MockInstall(module, db)

	// create namespace so we have something to update
	outcome, err := module.NamespaceCreate(context.Background(), input)
	ass.NoError(err)
	ass.Nil(outcome.Errors)

	t.Run("happy case", func(t *testing.T) {
		_, err = module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
			NamespaceID:      outcome.Namespace.ID,
			NamespaceVersion: outcome.Namespace.Version,
			Object: &dto.NamespaceUpdateInputObject{
				Features: &dto.NamespaceUpdateInputFeatures{
					Register: util.NilBool(true),
				},
			},
		})

		features, err := module.Features(context.Background(), outcome.Namespace)
		ass.NoError(err)
		ass.True(features.Register)
	})

	t.Run("update with invalid version -> conflict", func(t *testing.T) {
		_, err = module.NamespaceUpdate(context.Background(), dto.NamespaceUpdateInput{
			NamespaceID:      outcome.Namespace.ID,
			NamespaceVersion: "invalid-version",
			Object: &dto.NamespaceUpdateInputObject{
				Features: &dto.NamespaceUpdateInputFeatures{
					Register: util.NilBool(true),
				},
			},
		})

		ass.Equal(err, util.ErrorVersionConflict)
	})
}
