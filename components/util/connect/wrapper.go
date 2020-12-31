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

func NewWrapper(config map[string]DatabaseConfig) Wrapper {
	return Wrapper{
		config:      config,
		connections: &sync.Map{},
	}
}

type Wrapper struct {
	config      map[string]DatabaseConfig
	connections *sync.Map
}

func (dbs *Wrapper) Master() (*gorm.DB, error) {
	return dbs.get("master")
}

func (dbs *Wrapper) get(name string) (*gorm.DB, error) {
	if db, ok := dbs.connections.Load(name); ok {
		return db.(*gorm.DB), nil // Connection already established
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

func (dbs *Wrapper) dialector(cnf DatabaseConfig) gorm.Dialector {
	switch cnf.Driver {
	case "sqlite3":
		return sqlite.Open(cnf.Url)

	case "postgres":
		return postgres.Open(cnf.Url)

	default:
		panic("unsupported driver: " + cnf.Driver)
	}
}
