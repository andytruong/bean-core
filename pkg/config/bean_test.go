package config

import (
	"context"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

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

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	db := util.MockDatabase()
	util.MockInstall(this, db)

	err := connect.Transaction(
		context.Background(),
		db,
		func(tx *gorm.DB) error {
			hostId := this.id.MustULID()
			access := api.AccessMode("444")
			outcome, err := this.BucketCreate(context.Background(), tx, dto.BucketCreateInput{
				HostId:      hostId,
				Slug:        util.NilString("qa"),
				Title:       util.NilString("QA"),
				Description: util.NilString("Just for QA"),
				Access:      &access,
			})

			ass.NoError(err)
			ass.Empty(outcome.Errors)
			ass.Equal(hostId, outcome.Bucket.HostId)
			ass.Equal("qa", outcome.Bucket.Slug)
			ass.Equal("QA", outcome.Bucket.Title)
			ass.Equal("Just for QA", *outcome.Bucket.Description)
			ass.Equal(access, outcome.Bucket.Access)

			return err
		},
	)

	ass.NoError(err)
}

func Test_Update(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	db := util.MockDatabase()
	util.MockInstall(this, db)

	err := connect.Transaction(
		context.Background(),
		db,
		func(tx *gorm.DB) error {
			privateAccess := api.AccessModePrivate
			oCreate, _ := this.BucketCreate(context.Background(), tx, dto.BucketCreateInput{
				HostId:      this.id.MustULID(),
				Slug:        util.NilString("qa"),
				Title:       util.NilString("QA"),
				Description: util.NilString("Just for QA"),
				Access:      &privateAccess,
			})

			publicAccess := api.AccessModePublicRead
			outcome, err := this.BucketUpdate(context.Background(), tx, dto.BucketUpdateInput{
				Id:          oCreate.Bucket.Id,
				Version:     oCreate.Bucket.Version,
				Title:       util.NilString("Test"),
				Description: util.NilString("Just for Testing"),
				Access:      &publicAccess,
			})

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
}
