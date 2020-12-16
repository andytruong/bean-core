package infra

import (
	"sync"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	util2 "bean/components/util"
)

type databases struct {
	config      map[string]DatabaseConfig
	connections *sync.Map
}

func (this *databases) master() (*gorm.DB, error) {
	return this.get("master")
}

func (this *databases) get(name string) (*gorm.DB, error) {
	if db, ok := this.connections.Load(name); ok {
		// Connection already established
		return db.(*gorm.DB), nil
	} else if cnf, ok := this.config[name]; !ok {
		// No configuration found for requested-DB
		return nil, errors.Wrap(util2.ErrorConfig, "database config not provided: "+name)
	} else {
		if con, err := gorm.Open(this.dialector(cnf), &gorm.Config{SkipDefaultTransaction: true}); nil != err {
			return nil, err
		} else {
			this.connections.Store(name, con)

			return con, nil
		}
	}
}

func (this *databases) dialector(cnf DatabaseConfig) gorm.Dialector {
	switch cnf.Driver {
	case "sqlite3":
		return sqlite.Open(cnf.Url)

	case "postgres":
		return postgres.Open(cnf.Url)

	default:
		panic("unsupported driver: " + cnf.Driver)
	}
}
