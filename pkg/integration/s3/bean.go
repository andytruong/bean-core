package s3

import (
	"gorm.io/gorm"

	"bean/pkg/util"
)

type S3IntegrationBean struct {
}

func (this S3IntegrationBean) Migrate(tx *gorm.DB, driver string) error {
	panic("implement me")
}

func (this S3IntegrationBean) Dependencies() []util.Bean {
	panic("implement me")
}
