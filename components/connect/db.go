package connect

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"bean/components/util"
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

func NewWrapper(cnf map[string]DatabaseConfig) *Wrapper {
	return &Wrapper{
		cnf:         cnf,
		dbs:         &sync.Map{},
		prepareStmt: true,
	}
}

type Wrapper struct {
	cnf map[string]DatabaseConfig
	dbs *sync.Map

	prepareStmt bool
}

func (w *Wrapper) PrepareStmt(value bool) *Wrapper {
	w.prepareStmt = value

	return w
}

func (w *Wrapper) Master() (*gorm.DB, error) {
	return w.get("master")
}

func (w *Wrapper) get(name string) (*gorm.DB, error) {
	if db, ok := w.dbs.Load(name); ok {
		return db.(*gorm.DB), nil // Connection already established
	} else if cnf, ok := w.cnf[name]; !ok {
		// No configuration found for requested-DB
		return nil, errors.Wrap(util.ErrorConfig, "database config not provided: "+name)
	} else {
		dbCnf := &gorm.Config{
			SkipDefaultTransaction:                   true,
			DisableAutomaticPing:                     true,
			DisableForeignKeyConstraintWhenMigrating: true,
			PrepareStmt:                              w.prepareStmt,
		}

		if con, err := gorm.Open(w.dialector(cnf), dbCnf); nil != err {
			return nil, err
		} else {
			w.dbs.Store(name, con)

			return con, nil
		}
	}
}

func (w *Wrapper) dialector(cnf DatabaseConfig) gorm.Dialector {
	switch cnf.Driver {
	case "sqlite3":
		return sqlite.Open(cnf.Url)

	case "postgres":
		return postgres.Open(cnf.Url)

	default:
		panic("unsupported driver: " + cnf.Driver)
	}
}

func Transaction(ctx context.Context, callback func(tx *gorm.DB) error) error {
	con := DB(ctx)
	txn := con.Begin()
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
