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
