package infra

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"bean/pkg/util"
)


type databases struct {
	config      map[string]DatabaseConfig
	connections *sync.Map
}

func (this *databases) get(name string) (*gorm.DB, error) {
	if db, ok := this.connections.Load(name); ok {
		return db.(*gorm.DB), nil
	} else if cnf, ok := this.config[name]; !ok {
		return nil, errors.Wrap(util.ConfigError, "database config not provided: "+name)
	} else if con, err := gorm.Open(cnf.Driver, cnf.Url); nil != err {
		return nil, err
	} else {
		this.connections.Store(name, con)

		return con, nil
	}
}
