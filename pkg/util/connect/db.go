package connect

import (
	"context"

	"gorm.io/gorm"
)

const (
	// Driver names
	SQLite   = "sqlite3"
	Postgres = "postgres"

	// Table names
	TableConfigBucket             = "config_buckets"
	TableConfigVariable           = "config_variables"
	TableAccessSession            = "access_session"
	TableNamespace                = "namespaces"
	TableNamespaceMemberships     = "namespace_memberships"
	TableNamespaceDomains         = "namespace_domains"
	TableNamespaceConfig          = "namespace_config"
	TableManagerEdge              = "namespace_manager_edge"
	TableUserEmail                = "user_emails"
	TableAccessPassword           = "user_passwords"
	TableUserEmailUnverified      = "user_unverified_emails"
	TableIntegrationS3            = "s3_application"
	TableIntegrationS3Credentials = "s3_credentials"
)

func Transaction(ctx context.Context, db *gorm.DB, callback func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	er := callback(tx)

	if nil != er {
		return tx.Rollback().Error
	} else {
		return tx.Commit().Error
	}
}
