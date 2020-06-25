package s3

import (
	"testing"

	"bean/pkg/util"
)

func bean() *S3IntegrationBean {
	db := util.MockDatabase()
	id := util.MockIdentifier()
	logger := util.MockLogger()
	bean := NewS3Integration(db, id, logger)
	util.MockInstall(bean, db)

	return bean
}

func Test(t *testing.T) {
	this := bean()

	t.Run("DB schema", func(t *testing.T) {
		this.db.Migrator().HasTable("s3_application")
	})
}
