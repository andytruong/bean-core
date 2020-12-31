package connect

import (
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"bean/components/util"
)

func NewWrapper(cnf map[string]DatabaseConfig) *Wrapper {
	return &Wrapper{
		cnf: cnf,
		dbs: &sync.Map{},
	}
}

type Wrapper struct {
	cnf map[string]DatabaseConfig
	dbs *sync.Map
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
			PrepareStmt:                              true,
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
