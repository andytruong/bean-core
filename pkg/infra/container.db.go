package infra

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

type databases struct {
	config      map[string]DatabaseConfig
	connections *sync.Map
}

func (dbs *databases) master() (*gorm.DB, error) {
	return dbs.get("master")
}

func (dbs *databases) get(name string) (*gorm.DB, error) {
	if db, ok := dbs.connections.Load(name); ok {
		// Connection already established
		return db.(*gorm.DB), nil
	} else if cnf, ok := dbs.config[name]; !ok {
		// No configuration found for requested-DB
		return nil, errors.Wrap(util.ErrorConfig, "database config not provided: "+name)
	} else {
		if con, err := gorm.Open(dbs.dialector(cnf), &gorm.Config{SkipDefaultTransaction: true}); nil != err {
			return nil, err
		} else {
			dbs.connections.Store(name, con)

			return con, nil
		}
	}
}

func (dbs *databases) dialector(cnf DatabaseConfig) gorm.Dialector {
	switch cnf.Driver {
	case "sqlite3":
		return sqlite.Open(cnf.Url)

	case "postgres":
		return postgres.Open(cnf.Url)

	default:
		panic("unsupported driver: " + cnf.Driver)
	}
}
