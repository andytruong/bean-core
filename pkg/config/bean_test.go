package config

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

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
				outcome, err := this.Bucket.Create(ctx, tx, dto.BucketCreateInput{
					HostId:      hostId,
					Slug:        util.NilString("doe"),
					Title:       util.NilString("Doe"),
					Description: util.NilString("Just for John Doe"),
					Access:      &access,
					Schema:      `{"type:"number"}`,
				})

				ass.NoError(err)
				ass.Empty(outcome.Errors)
				ass.Equal(hostId, outcome.Bucket.HostId)
				ass.Equal("doe", outcome.Bucket.Slug)
				ass.Equal("Doe", outcome.Bucket.Title)
				ass.Equal("Just for John Doe", *outcome.Bucket.Description)
				ass.Equal(access, outcome.Bucket.Access)

				return err
			},
		)

		ass.NoError(err)
	})

	t.Run("bucket.update", func(t *testing.T) {
		err := connect.Transaction(
			ctx,
			db,
			func(tx *gorm.DB) error {
				privateAccess := api.AccessModePrivate
				oCreate, _ := this.Bucket.Create(ctx, tx, dto.BucketCreateInput{
					HostId:      this.id.MustULID(),
					Slug:        util.NilString("qa"),
					Title:       util.NilString("QA"),
					Description: util.NilString("Just for QA"),
					Access:      &privateAccess,
					Schema:      `{"type:"number"}`,
				})

				publicAccess := api.AccessModePublicRead
				outcome, err := this.Bucket.Update(ctx, tx, dto.BucketUpdateInput{
					Id:          oCreate.Bucket.Id,
					Version:     oCreate.Bucket.Version,
					Title:       util.NilString("Test"),
					Description: util.NilString("Just for Testing"),
					Access:      &publicAccess,
				})

				ass.NotNil(outcome)
				ass.Empty(outcome.Errors)
				ass.NotEqual(oCreate.Bucket.Version, outcome.Bucket.Version)
				ass.Equal(oCreate.Bucket.Slug, outcome.Bucket.Slug)
				ass.Equal("Test", outcome.Bucket.Title)
				ass.Equal("Just for Testing", *outcome.Bucket.Description)
				ass.Equal(publicAccess, outcome.Bucket.Access)

				return err
			},
		)

		ass.NoError(err)
	})

	t.Run("bucket.load", func(t *testing.T) {
		var err error
		var oCreate *dto.BucketMutationOutcome
		var bucket *model.ConfigBucket
		hostId := this.id.MustULID()
		access := api.AccessMode("444")
		tx := db.BeginTx(ctx, &sql.TxOptions{})
		defer tx.Rollback()

		oCreate, err = this.Bucket.Create(ctx, tx, dto.BucketCreateInput{
			HostId:      hostId,
			Slug:        util.NilString("load-doe"),
			Title:       util.NilString("Doe"),
			Description: util.NilString("Just for John Doe"),
			Access:      &access,
			Schema:      `{"type:"number"}`,
		})

		ass.NoError(err)
		ass.NotNil(oCreate)

		bucket, err = this.Bucket.BucketLoad(context.Background(), tx, oCreate.Bucket.Id)
		ass.NoError(err)
		ass.Equal(hostId, bucket.HostId)
		ass.Equal("load-doe", bucket.Slug)
		ass.Equal("Doe", bucket.Title)
		ass.Equal("Just for John Doe", *bucket.Description)
		ass.Equal(access, bucket.Access)
	})
}

func Test_Variable(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	db := util.MockDatabase()
	util.MockInstall(this, db)

	t.Run("variable.create", func(t *testing.T) {
		t.Run("on read-only bucket", func(t *testing.T) {
			// create read-only bucket
		})

		t.Run("on writable bucket", func(t *testing.T) {
			// create writable bucket
		})
	})

	t.Run("variable.load", func(t *testing.T) {
		t.Run("load on private bucket", func(t *testing.T) {
			// create private bucket
			// create variable
			// load & assert outcome
		})

		t.Run("load on read only bucket", func(t *testing.T) {
			// create read-only bucket
			// create variable
			// load & assert outcome
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
