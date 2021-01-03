package config

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"bean/components/claim"
	"bean/components/connect"
	"bean/components/scalar"
	"bean/components/util"
	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
)

func bundle() *ConfigBundle {
	idr := util.MockIdentifier()
	log := util.MockLogger()
	bun := NewConfigBundle(idr, log)

	return bun
}

func Test_Bucket(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	t.Run("create", func(t *testing.T) {
		t.Run("invalid schema", func(t *testing.T) {
			hostId := bundle.idr.MustULID()
			access := scalar.AccessMode("444")
			out, err := bundle.BucketService.Create(ctx, dto.BucketCreateInput{
				HostId:      hostId,
				Slug:        scalar.NilString("doe"),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type":invalid}`,
			})

			ass.NoError(err)
			ass.NotNil(out)
			ass.Len(out.Errors, 1)
			ass.Contains(out.Errors[0].Message, "invalid character 'i' looking for beginning of value")
		})

		t.Run("valid schema", func(t *testing.T) {
			hostId := bundle.idr.MustULID()
			access := scalar.AccessMode("444")
			out, err := bundle.BucketService.Create(ctx, dto.BucketCreateInput{
				HostId:      hostId,
				Slug:        scalar.NilString("doe"),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type":"number"}`,
			})

			ass.NoError(err)
			ass.Empty(out.Errors)
			ass.Equal(hostId, out.Bucket.HostId)
			ass.Equal("doe", out.Bucket.Slug)
			ass.Equal("Doe", out.Bucket.Title)
			ass.Equal("Just for John Doe", *out.Bucket.Description)
			ass.Equal(access, out.Bucket.Access)
		})
	})

	t.Run("update", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		privateAccess := scalar.AccessModePrivate
		oCreate, _ := bundle.BucketService.Create(connect.DBToContext(ctx, tx), dto.BucketCreateInput{
			HostId:      bundle.idr.MustULID(),
			Slug:        scalar.NilString("qa"),
			Title:       scalar.NilString("QA"),
			Description: scalar.NilString("Just for QA"),
			Access:      &privateAccess,
			Schema:      `{"type":"number"}`,
			IsPublished: false,
		})

		publicAccess := scalar.AccessModePublicRead
		oUpdate, err := bundle.BucketService.Update(connect.DBToContext(ctx, tx), dto.BucketUpdateInput{
			Id:          oCreate.Bucket.Id,
			Version:     oCreate.Bucket.Version,
			Title:       scalar.NilString("Test"),
			Description: scalar.NilString("Just for Testing"),
			Access:      &publicAccess,
			Schema:      scalar.NilString(`{"type":"string"}`),
			IsPublished: scalar.NilBool(true),
		})

		ass.NoError(err)
		ass.NotNil(oUpdate)
		ass.Empty(oUpdate.Errors)
		ass.NotEqual(oCreate.Bucket.Version, oUpdate.Bucket.Version)
		ass.Equal(oCreate.Bucket.Slug, oUpdate.Bucket.Slug)
		ass.Equal("Test", oUpdate.Bucket.Title)
		ass.Equal("Just for Testing", *oUpdate.Bucket.Description)
		ass.Equal(publicAccess, oUpdate.Bucket.Access)

		t.Run("can't unpublished a published bucket", func(t *testing.T) {
			_, err := bundle.BucketService.Update(connect.DBToContext(ctx, tx), dto.BucketUpdateInput{
				Id:          oUpdate.Bucket.Id,
				Version:     oUpdate.Bucket.Version,
				IsPublished: scalar.NilBool(false),
			})

			ass.Error(err)
			ass.Equal(err.Error(), "change not un-publish a published bucket: locked")
		})

		t.Run("can't change schema is isPublished on", func(t *testing.T) {
			_, err := bundle.BucketService.Update(connect.DBToContext(ctx, tx), dto.BucketUpdateInput{
				Id:          oCreate.Bucket.Id,
				Version:     oCreate.Bucket.Version,
				Title:       scalar.NilString("Test"),
				Description: scalar.NilString("Just for Testing"),
				Access:      &publicAccess,
				Schema:      scalar.NilString(`{"type":"int"}`),
			})

			ass.Error(err)
			ass.Equal(util.ErrorVersionConflict.Error(), err.Error())
		})
	})

	t.Run("load", func(t *testing.T) {
		var err error
		var oCreate *dto.BucketMutationOutcome
		var bucket *model.ConfigBucket
		hostId := bundle.idr.MustULID()
		access := scalar.AccessMode("444")
		tx := db.Begin()
		defer tx.Rollback()

		oCreate, err = bundle.BucketService.Create(connect.DBToContext(ctx, tx), dto.BucketCreateInput{
			HostId:      hostId,
			Slug:        scalar.NilString("load-doe"),
			Title:       scalar.NilString("Doe"),
			Description: scalar.NilString("Just for John Doe"),
			Access:      &access,
			Schema:      `{"type":"number"}`,
			IsPublished: true,
		})

		ass.NoError(err)
		ass.NotNil(oCreate)

		t.Run("load by ID", func(t *testing.T) {
			bucket, err = bundle.BucketService.Load(connect.DBToContext(context.Background(), tx), dto.BucketKey{Id: oCreate.Bucket.Id})
			ass.NoError(err)
			ass.Equal(hostId, bucket.HostId)
			ass.Equal("load-doe", bucket.Slug)
			ass.Equal("Doe", bucket.Title)
			ass.Equal("Just for John Doe", *bucket.Description)
			ass.Equal(access, bucket.Access)
			ass.Equal(true, bucket.IsPublished)
		})

		t.Run("load by slug", func(t *testing.T) {
			bucket, err = bundle.BucketService.Load(connect.DBToContext(context.Background(), tx), dto.BucketKey{Slug: oCreate.Bucket.Slug})
			ass.NoError(err)
			ass.Equal(hostId, bucket.HostId)
			ass.Equal("load-doe", bucket.Slug)
			ass.Equal("Doe", bucket.Title)
			ass.Equal("Just for John Doe", *bucket.Description)
			ass.Equal(access, bucket.Access)
			ass.Equal(true, bucket.IsPublished)
		})
	})
}

func Test_Variable(t *testing.T) {
	ass := assert.New(t)
	bundle := bundle()
	db := connect.MockDatabase()
	ctx := connect.DBToContext(context.Background(), db)
	connect.MockInstall(ctx, bundle)

	t.Run("create", func(t *testing.T) {
		t.Run("read-only bucket", func(t *testing.T) {
			ctx := context.Background()
			tx := db.Begin(&sql.TxOptions{})
			defer tx.Rollback()
			hostId := bundle.idr.MustULID()
			access := scalar.AccessModePrivateReadonly

			// create read-only bucket
			oCreate, err := bundle.BucketService.Create(connect.DBToContext(ctx, tx), dto.BucketCreateInput{
				HostId:      hostId,
				Slug:        scalar.NilString("load-doe"),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type":"number"}`,
				IsPublished: true,
			})

			ass.NoError(err, "no err on create read-only bucket")
			ass.Empty(oCreate.Errors)

			// create variable
			out, err := bundle.VariableService.Create(connect.DBToContext(ctx, tx), dto.VariableCreateInput{
				BucketId:    oCreate.Bucket.Id,
				Name:        "foo",
				Description: nil,
				Value:       "1",
				IsLocked:    scalar.NilBool(false),
			})

			// assert error
			ass.Error(err)
			ass.Equal(util.ErrorAccessDenied, err)
			ass.Nil(out)
		})

		t.Run("writable bucket", func(t *testing.T) {
			// setup auth context
			userId := bundle.idr.MustULID()
			claims := claim.NewPayload()
			claims.SetUserId(userId)
			ctx := claim.PayloadToContext(context.Background(), &claims)

			tx := db.Begin()
			defer tx.Rollback()
			access := scalar.AccessModePrivate

			// create read-only bucket
			oCreate, err := bundle.BucketService.Create(connect.DBToContext(ctx, tx), dto.BucketCreateInput{
				HostId:      userId,
				Slug:        scalar.NilString("load-doe"),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type":"number"}`,
				IsPublished: true,
			})

			ass.NoError(err, "no err on create read-only bucket")
			ass.Empty(oCreate.Errors)

			t.Run("invalid schema", func(t *testing.T) {
				// create variable
				out, err := bundle.VariableService.Create(connect.DBToContext(ctx, tx), dto.VariableCreateInput{
					BucketId:    oCreate.Bucket.Id,
					Name:        "foo",
					Description: nil,
					Value:       `"Should be a number"`,
					IsLocked:    scalar.NilBool(false),
				})

				// assert error
				ass.NoError(err)
				ass.NotNil(out)
				ass.Nil(out.Variable)
				ass.Len(out.Errors, 1)
				ass.Contains(out.Errors[0].Message, `"Should be a number" type should be number, got string`)
			})

			t.Run("valid schema", func(t *testing.T) {
				// create variable
				out, err := bundle.VariableService.Create(connect.DBToContext(ctx, tx), dto.VariableCreateInput{
					BucketId:    oCreate.Bucket.Id,
					Name:        "foo",
					Description: nil,
					Value:       "1",
					IsLocked:    scalar.NilBool(false),
				})

				// assert error
				ass.NoError(err)
				ass.NotNil(out)
				ass.Empty(out.Errors)
				ass.Equal("1", out.Variable.Value)
			})
		})
	})

	t.Run("load", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		setup := func(access scalar.AccessMode) (context.Context, *model.ConfigBucket, *model.ConfigVariable) {
			// setup auth context
			authorId := bundle.idr.MustULID()
			authorClaims := claim.NewPayload()
			authorClaims.SetUserId(authorId)
			authorCtx := claim.PayloadToContext(context.Background(), &authorClaims)

			// create private bucket
			oBucketCreate, err := bundle.BucketService.Create(connect.DBToContext(authorCtx, tx), dto.BucketCreateInput{
				HostId:      authorId,
				Slug:        scalar.NilString(bundle.idr.MustULID()),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type":"number"}`,
				IsPublished: true,
			})

			ass.NoError(err)

			// create variable
			oVarCreate, err := bundle.VariableService.Create(connect.DBToContext(authorCtx, tx), dto.VariableCreateInput{
				BucketId:    oBucketCreate.Bucket.Id,
				Name:        "foo",
				Description: nil,
				Value:       "1",
				IsLocked:    scalar.NilBool(false),
			})

			ass.NoError(err)

			return authorCtx, oBucketCreate.Bucket, oVarCreate.Variable
		}

		t.Run("private bucket", func(t *testing.T) {
			_, _, variable := setup(scalar.AccessModePrivate)

			// load & assert outcome
			load, err := bundle.VariableService.Load(ctx, dto.VariableKey{Id: variable.Id})
			ass.Error(err)
			ass.Nil(load)
		})

		t.Run("read-only bucket", func(t *testing.T) {
			ctx, bucket, variable := setup(scalar.AccessModePrivate)

			t.Run("by ID", func(t *testing.T) {
				load, err := bundle.VariableService.Load(connect.DBToContext(ctx, tx), dto.VariableKey{Id: variable.Id})
				ass.NoError(err)
				ass.Equal(bucket.Id, load.BucketId)
				ass.Equal("1", load.Value)
			})

			t.Run("by name", func(t *testing.T) {
				load, err := bundle.VariableService.Load(connect.DBToContext(ctx, tx), dto.VariableKey{BucketId: bucket.Id, Name: variable.Name})
				ass.NoError(err)
				ass.Equal(bucket.Id, load.BucketId)
				ass.Equal("1", load.Value)
			})
		})
	})

	t.Run("update", func(t *testing.T) {
		t.Run("read-only bucket", func(t *testing.T) {
			// TODO
			// create read-only bucket
			// create variable
			// update variable & assert outcome
		})

		t.Run("writable bucket", func(t *testing.T) {
			// TODO
			// create read-only bucket
			// create variable
			// update variable & assert outcome
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("delete on read-only bucket", func(t *testing.T) {
			// TODO
			// create read-only bucket
			// create variable
			// delete variable & assert outcome
		})

		t.Run("delete on writable bucket", func(t *testing.T) {
			t.Run("write on locked variable", func(t *testing.T) {
				// TODO
				// create read-only bucket
				// create variable
				// delete variable & assert outcome
			})

			t.Run("write on unlocked variable", func(t *testing.T) {
				// TODO
				// create read-only bucket
				// create variable
				// delete variable & assert outcome
			})
		})
	})
}
