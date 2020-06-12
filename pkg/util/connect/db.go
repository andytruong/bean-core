package connect

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
)

const (
	// Driver names
	SQLite   = "sqlite3"
	Postgres = "postgres"

	// Table names
	TableConfigBucket         = "config_buckets"
	TableConfigVariable       = "config_variables"
	TableAccessSession        = "access_session"
	TableAccessPassword       = "user_password"
	TableNamespace            = "namespaces"
	TableNamespaceMemberships = "namespace_memberships"
	TableNamespaceDomains     = "namespace_domains"
	TableNamespaceConfig      = "namespace_config"
	TableManagerEdge          = "namespace_manager_edge"
	TableUserEmail            = "user_emails"
	TableUserEmailUnverified  = "user_unverified_emails"
)

func Transaction(ctx context.Context, db *gorm.DB, callback func(tx *gorm.DB) error) error {
	tx := db.BeginTx(ctx, &sql.TxOptions{})
	er := callback(tx)

	if nil != er {
		return tx.Rollback().Error
	} else {
		return tx.Commit().Error
	}
}
