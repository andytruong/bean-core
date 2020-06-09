package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

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

func Test_Create(t *testing.T) {
	ass := assert.New(t)
	this := bean()
	db := util.MockDatabase().LogMode(false)
	util.MockInstall(this, db)

	err := connect.Transaction(
		context.Background(),
		db,
		func(tx *gorm.DB) error {
			outcome, err := this.BucketCreate(context.Background(), tx, dto.BucketCreateInput{
				HostId:      this.id.MustULID(),
				Slug:        nil,
				Title:       nil,
				Description: nil,
				Access:      nil,
			})

			ass.NoError(err)
			fmt.Println("outcome: ", outcome)

			return err
		},
	)

	ass.NoError(err)
}
