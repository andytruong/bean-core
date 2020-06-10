package config

import (
	"context"
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
				outcome, err := this.Bucket.Create(context.Background(), tx, dto.BucketCreateInput{
					HostId:      hostId,
					Slug:        util.NilString("doe"),
					Title:       util.NilString("Doe"),
					Description: util.NilString("Just for John Doe"),
					Access:      &access,
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
			context.Background(),
			db,
			func(tx *gorm.DB) error {
				privateAccess := api.AccessModePrivate
				oCreate, _ := this.Bucket.Create(context.Background(), tx, dto.BucketCreateInput{
					HostId:      this.id.MustULID(),
					Slug:        util.NilString("qa"),
					Title:       util.NilString("QA"),
					Description: util.NilString("Just for QA"),
					Access:      &privateAccess,
				})

				publicAccess := api.AccessModePublicRead
				outcome, err := this.Bucket.Update(context.Background(), tx, dto.BucketUpdateInput{
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

		_ = connect.Transaction(
			context.Background(),
			db,
			func(tx *gorm.DB) error {

				oCreate, err = this.Bucket.Create(context.Background(), tx, dto.BucketCreateInput{
					HostId:      hostId,
					Slug:        util.NilString("doe"),
					Title:       util.NilString("Doe"),
					Description: util.NilString("Just for John Doe"),
					Access:      &access,
				})

				return err
			},
		)

		bucket, err = this.Bucket.BucketLoad(context.Background(), db, oCreate.Bucket.Id)
		ass.NoError(err)
		ass.Equal(hostId, bucket.HostId)
		ass.Equal("doe", bucket.Slug)
		ass.Equal("Doe", bucket.Title)
		ass.Equal("Just for John Doe", *bucket.Description)
		ass.Equal(access, bucket.Access)
	})
}
