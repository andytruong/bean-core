package s3

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bean/pkg/util"
)

func bean() *S3IntegrationBean {
	db := util.MockDatabase()
	id := util.MockIdentifier()
	logger := util.MockLogger()
	bean := NewS3Integration(db, id, logger, &Genetic{
		Key: []byte("01EBWB516AP6BQD7"),
	})
	util.MockInstall(bean, db)

	return bean
}

func Test(t *testing.T) {
	ass := assert.New(t)
	this := bean()

	t.Run("DB schema", func(t *testing.T) {
		this.db.Migrator().HasTable("s3_application")
	})

	t.Run("core app", func(t *testing.T) {
		t.Run("encrypt", func(t *testing.T) {
			encrypted := this.coreApp.encrypt("xxxxxxxxxxxxxxxxxxxxx")
			decrypted := this.coreApp.decrypt(encrypted)

			ass.Equal("xxxxxxxxxxxxxxxxxxxxx", decrypted)
			ass.True(len(encrypted)*2 <= 256)
		})
	})
}
