package config

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/api"
	"bean/pkg/util/connect"
)

func bean() *ConfigBean {
	id := util.MockIdentifier()
	logger := util.MockLogger()
	bean := NewConfigBean(id, logger)

	return bean
}

func Test_Bucket(t *testing.T) {
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
				access := api.AccessMode("444")
				out, err := this.CoreBucket.Create(ctx, tx, dto.BucketCreateInput{
					HostId:      hostId,
					Slug:        util.NilString("doe"),
					Title:       util.NilString("Doe"),
					Description: util.NilString("Just for John Doe"),
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

		privateAccess := api.AccessModePrivate
		oCreate, _ := this.CoreBucket.Create(ctx, tx, dto.BucketCreateInput{
			HostId:      this.id.MustULID(),
			Slug:        util.NilString("qa"),
			Title:       util.NilString("QA"),
			Description: util.NilString("Just for QA"),
			Access:      &privateAccess,
			Schema:      `{"type:"number"}`,
			IsPublished: false,
		})

		publicAccess := api.AccessModePublicRead
		oUpdate, err := this.CoreBucket.Update(ctx, tx, dto.BucketUpdateInput{
			Id:          oCreate.Bucket.Id,
			Version:     oCreate.Bucket.Version,
			Title:       util.NilString("Test"),
			Description: util.NilString("Just for Testing"),
			Access:      &publicAccess,
			Schema:      util.NilString(`{"type":"string"}`),
			IsPublished: util.NilBool(true),
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
				IsPublished: util.NilBool(false),
			})

			ass.Error(err)
			ass.Equal(err.Error(), "change not un-publish a published bucket: locked")
		})

		t.Run("can't change schema is isPublished on", func(t *testing.T) {
			_, err := this.CoreBucket.Update(ctx, tx, dto.BucketUpdateInput{
				Id:          oCreate.Bucket.Id,
				Version:     oCreate.Bucket.Version,
				Title:       util.NilString("Test"),
				Description: util.NilString("Just for Testing"),
				Access:      &publicAccess,
				Schema:      util.NilString(`{"type":"int"}`),
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
		access := api.AccessMode("444")
		tx := db.Begin()
		defer tx.Rollback()

		oCreate, err = this.CoreBucket.Create(ctx, tx, dto.BucketCreateInput{
			HostId:      hostId,
			Slug:        util.NilString("load-doe"),
			Title:       util.NilString("Doe"),
			Description: util.NilString("Just for John Doe"),
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
			access := api.AccessModePrivateReadonly

			// create read-only bucket
			oCreate, err := this.CoreBucket.Create(ctx, tx, dto.BucketCreateInput{
				HostId:      hostId,
				Slug:        util.NilString("load-doe"),
				Title:       util.NilString("Doe"),
				Description: util.NilString("Just for John Doe"),
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
				IsLocked:    util.NilBool(false),
			})

			// assert error
			ass.Error(err)
			ass.Equal(util.ErrorAccessDenied, err)
			ass.Nil(out)
		})

		t.Run("on writable bucket", func(t *testing.T) {
			userId := this.id.MustULID()
			ctx := context.WithValue(context.Background(), util.CxtKeyClaims, &util.Claims{
				StandardClaims: jwt.StandardClaims{Subject: userId},
			})
			tx := db.Begin()
			defer tx.Rollback()
			access := api.AccessModePrivate

			// create read-only bucket
			oCreate, err := this.CoreBucket.Create(ctx, tx, dto.BucketCreateInput{
				HostId:      userId,
				Slug:        util.NilString("load-doe"),
				Title:       util.NilString("Doe"),
				Description: util.NilString("Just for John Doe"),
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
				IsLocked:    util.NilBool(false),
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

		setup := func(access api.AccessMode) (context.Context, *model.ConfigBucket, *model.ConfigVariable) {
			authorId := this.id.MustULID()
			authorClaims := &util.Claims{}
			authorClaims.Subject = authorId
			authorCtx := context.WithValue(context.Background(), util.CxtKeyClaims, authorClaims)

			// create private bucket
			oBucketCreate, err := this.CoreBucket.Create(authorCtx, tx, dto.BucketCreateInput{
				HostId:      authorId,
				Slug:        util.NilString(this.id.MustULID()),
				Title:       util.NilString("Doe"),
				Description: util.NilString("Just for John Doe"),
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
				IsLocked:    util.NilBool(false),
			})

			ass.NoError(err)

			return authorCtx, oBucketCreate.Bucket, oVarCreate.Variable
		}

		t.Run("load on private bucket", func(t *testing.T) {
			_, _, variable := setup(api.AccessModePrivate)

			// load & assert outcome
			otherCtx := context.Background()
			load, err := this.CoreVariable.Load(otherCtx, tx, variable.Id)
			ass.Error(err)
			ass.Nil(load)
		})

		t.Run("load on read only bucket", func(t *testing.T) {
			ctx, bucket, variable := setup(api.AccessModePrivate)

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
