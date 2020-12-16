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
	TableSpace               = "spaces"
	TableSpaceMemberships    = "space_memberships"
	TableUserEmail           = "user_emails"
	TableUserEmailUnverified = "user_unverified_emails"
)

func Transaction(ctx context.Context, db *gorm.DB, callback func(tx *gorm.DB) error) error {
	txn := db.WithContext(ctx).Begin()
	err := callback(txn)

	if nil != err {
		rollbackErr := txn.Rollback().Error
		if nil != rollbackErr {
			return rollbackErr
		}

		return err
	} else {
		return txn.Commit().Error
	}
}
