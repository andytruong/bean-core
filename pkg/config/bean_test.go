package config

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"bean/components/claim"
	"bean/components/scalar"
	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

func bean() *ConfigBean {
	id := util.MockIdentifier()
	logger := util.MockLogger()
	bean := NewConfigBean(id, logger)

	return bean
}

func Test_Bucket(t *testing.T) {
	t.Parallel()

	ass := assert.New(t)
	ctx := context.Background()
	this := bean()
	db := util.MockDatabase()
	util.MockInstall(this, db)

	t.Run("bucket.create", func(t *testing.T) {
		err := connect.Transaction(
			context.Background(),
			db,
			func(tx *gorm.DB) error {
				hostId := this.id.MustULID()
				access := scalar.AccessMode("444")
				out, err := this.CoreBucket.Create(tx, dto.BucketCreateInput{
					HostId:      hostId,
					Slug:        scalar.NilString("doe"),
					Title:       scalar.NilString("Doe"),
					Description: scalar.NilString("Just for John Doe"),
					Access:      &access,
					Schema:      `{"type:"number"}`,
				})

				ass.NoError(err)
				ass.Empty(out.Errors)
				ass.Equal(hostId, out.Bucket.HostId)
				ass.Equal("doe", out.Bucket.Slug)
				ass.Equal("Doe", out.Bucket.Title)
				ass.Equal("Just for John Doe", *out.Bucket.Description)
				ass.Equal(access, out.Bucket.Access)

				return err
			},
		)

		ass.NoError(err)
	})

	t.Run("bucket.update", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		privateAccess := scalar.AccessModePrivate
		oCreate, _ := this.CoreBucket.Create(tx, dto.BucketCreateInput{
			HostId:      this.id.MustULID(),
			Slug:        scalar.NilString("qa"),
			Title:       scalar.NilString("QA"),
			Description: scalar.NilString("Just for QA"),
			Access:      &privateAccess,
			Schema:      `{"type:"number"}`,
			IsPublished: false,
		})

		publicAccess := scalar.AccessModePublicRead
		oUpdate, err := this.CoreBucket.Update(ctx, tx, dto.BucketUpdateInput{
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
			_, err := this.CoreBucket.Update(ctx, tx, dto.BucketUpdateInput{
				Id:          oUpdate.Bucket.Id,
				Version:     oUpdate.Bucket.Version,
				IsPublished: scalar.NilBool(false),
			})

			ass.Error(err)
			ass.Equal(err.Error(), "change not un-publish a published bucket: locked")
		})

		t.Run("can't change schema is isPublished on", func(t *testing.T) {
			_, err := this.CoreBucket.Update(ctx, tx, dto.BucketUpdateInput{
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

	t.Run("bucket.load", func(t *testing.T) {
		var err error
		var oCreate *dto.BucketMutationOutcome
		var bucket *model.ConfigBucket
		hostId := this.id.MustULID()
		access := scalar.AccessMode("444")
		tx := db.Begin()
		defer tx.Rollback()

		oCreate, err = this.CoreBucket.Create(tx, dto.BucketCreateInput{
			HostId:      hostId,
			Slug:        scalar.NilString("load-doe"),
			Title:       scalar.NilString("Doe"),
			Description: scalar.NilString("Just for John Doe"),
			Access:      &access,
			Schema:      `{"type:"number"}`,
			IsPublished: true,
		})

		ass.NoError(err)
		ass.NotNil(oCreate)

		bucket, err = this.CoreBucket.Load(context.Background(), tx, oCreate.Bucket.Id)
		ass.NoError(err)
		ass.Equal(hostId, bucket.HostId)
		ass.Equal("load-doe", bucket.Slug)
		ass.Equal("Doe", bucket.Title)
		ass.Equal("Just for John Doe", *bucket.Description)
		ass.Equal(access, bucket.Access)
		ass.Equal(true, bucket.IsPublished)
	})
}

func Test_Variable(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	db := util.MockDatabase()
	util.MockInstall(this, db)

	t.Run("variable.create", func(t *testing.T) {
		t.Run("on read-only bucket", func(t *testing.T) {
			ctx := context.Background()
			tx := db.Begin(&sql.TxOptions{})
			defer tx.Rollback()
			hostId := this.id.MustULID()
			access := scalar.AccessModePrivateReadonly

			// create read-only bucket
			oCreate, err := this.CoreBucket.Create(tx, dto.BucketCreateInput{
				HostId:      hostId,
				Slug:        scalar.NilString("load-doe"),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type:"number"}`,
				IsPublished: true,
			})

			ass.NoError(err, "no err on create read-only bucket")
			ass.Empty(oCreate.Errors)

			// create variable
			out, err := this.CoreVariable.Create(ctx, tx, dto.VariableCreateInput{
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

		t.Run("on writable bucket", func(t *testing.T) {
			userId := this.id.MustULID()
			ctx := context.WithValue(context.Background(), claim.ContextKey, &claim.Payload{
				StandardClaims: jwt.StandardClaims{Subject: userId},
			})
			tx := db.Begin()
			defer tx.Rollback()
			access := scalar.AccessModePrivate

			// create read-only bucket
			oCreate, err := this.CoreBucket.Create(tx, dto.BucketCreateInput{
				HostId:      userId,
				Slug:        scalar.NilString("load-doe"),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type:"number"}`,
				IsPublished: true,
			})

			ass.NoError(err, "no err on create read-only bucket")
			ass.Empty(oCreate.Errors)

			// create variable
			out, err := this.CoreVariable.Create(ctx, tx, dto.VariableCreateInput{
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

	t.Run("variable.load", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		setup := func(access scalar.AccessMode) (context.Context, *model.ConfigBucket, *model.ConfigVariable) {
			authorId := this.id.MustULID()
			authorClaims := &claim.Payload{}
			authorClaims.Subject = authorId
			authorCtx := context.WithValue(context.Background(), claim.ContextKey, authorClaims)

			// create private bucket
			oBucketCreate, err := this.CoreBucket.Create(tx, dto.BucketCreateInput{
				HostId:      authorId,
				Slug:        scalar.NilString(this.id.MustULID()),
				Title:       scalar.NilString("Doe"),
				Description: scalar.NilString("Just for John Doe"),
				Access:      &access,
				Schema:      `{"type:"number"}`,
				IsPublished: true,
			})

			ass.NoError(err)

			// create variable
			oVarCreate, err := this.CoreVariable.Create(authorCtx, tx, dto.VariableCreateInput{
				BucketId:    oBucketCreate.Bucket.Id,
				Name:        "foo",
				Description: nil,
				Value:       "1",
				IsLocked:    scalar.NilBool(false),
			})

			ass.NoError(err)

			return authorCtx, oBucketCreate.Bucket, oVarCreate.Variable
		}

		t.Run("load on private bucket", func(t *testing.T) {
			_, _, variable := setup(scalar.AccessModePrivate)

			// load & assert outcome
			otherCtx := context.Background()
			load, err := this.CoreVariable.Load(otherCtx, tx, variable.Id)
			ass.Error(err)
			ass.Nil(load)
		})

		t.Run("load on read only bucket", func(t *testing.T) {
			ctx, bucket, variable := setup(scalar.AccessModePrivate)

			// load & assert outcome
			load, err := this.CoreVariable.Load(ctx, tx, variable.Id)
			ass.NoError(err)
			ass.Equal(bucket.Id, load.BucketId)
			ass.Equal("1", load.Value)
		})
	})

	t.Run("variable.update", func(t *testing.T) {
		t.Run("update on read-only bucket", func(t *testing.T) {
			// create read-only bucket
			// create variable
			// update variable & assert outcome
		})

		t.Run("update on writable bucket", func(t *testing.T) {
			// create read-only bucket
			// create variable
			// update variable & assert outcome
		})
	})

	t.Run("variable delete", func(t *testing.T) {
		t.Run("delete on read-only bucket", func(t *testing.T) {
			// create read-only bucket
			// create variable
			// delete variable & assert outcome
		})

		t.Run("delete on writable bucket", func(t *testing.T) {
			t.Run("write on locked variable", func(t *testing.T) {
				// create read-only bucket
				// create variable
				// delete variable & assert outcome
			})

			t.Run("write on unlocked variable", func(t *testing.T) {
				// create read-only bucket
				// create variable
				// delete variable & assert outcome
			})
		})
	})

	ass.True(true, "WIP")
}
