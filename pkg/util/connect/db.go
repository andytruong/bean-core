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
)

func Transaction(ctx context.Context, db *gorm.DB, callback func(tx *gorm.DB) error) error {
	txn := db.BeginTx(ctx, &sql.TxOptions{})
	err := callback(txn)

	if nil != err {
		return txn.Rollback().Error
	} else {
		return txn.Commit().Error
	}
}
