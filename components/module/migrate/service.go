package migrate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Runner struct {
	Tx     *gorm.DB
	Logger *zap.Logger
	Driver string
	Bean   string
	Dir    string
}

func (this Runner) Run() error {
	return filepath.Walk(
		this.Dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.HasSuffix(path, ".sql") {
					return this.installFile(path)
				}
			}
			
			return nil
		},
	)
}

func (this Runner) installFile(file string) error {
	migration := NewMigration(this.Bean, file)
	path := migration.RealPath()
	
	if !migration.DriverMatch(this.Driver) {
		this.Logger.Debug(
			"👉 driver unmatched",
			zap.String("bean", migration.Bean),
			zap.String("path", path),
		)
		
		return nil
	}
	
	if can, err := migration.IsExecuted(this.Tx); nil != err {
		return err
	} else if can {
		if migration.DriverMatch(this.Driver) {
			content, err := ioutil.ReadFile(path)
			if nil != err {
				return err
			}

			if err := this.Tx.Exec(string(content)).Error; nil != err {
				this.Logger.Info(
					"⚡️ executed migration",
					zap.String("bean", migration.Bean),
					zap.String("path", path),
				)
				
				return err
			} else {
				return migration.Save(this.Tx)
			}
		}
	}
	
	return nil
}
