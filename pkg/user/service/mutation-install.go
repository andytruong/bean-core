package service

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/util"
)

func NewUserInstallAPI(logger *zap.Logger) *UserInstallAPI {
	return &UserInstallAPI{logger: logger}
}

type UserInstallAPI struct {
	logger *zap.Logger
}

func (this UserInstallAPI) Run(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Walk(
			path.Dir(path.Dir(filename))+"/migration/",
			func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					if strings.HasSuffix(path, ".sql") {
						return this.install(tx, driver, path)
					}
				}

				return nil
			},
		)
	}

	return nil
}

func (this UserInstallAPI) install(tx *gorm.DB, driver string, file string) error {
	migration := util.NewMigration("user", file)
	path := migration.RealPath()

	if !migration.DriverMatch(driver) {
		this.logger.Info(
			"👉 driver unmatched",
			zap.String("module", migration.Module),
			zap.String("path", path),
		)

		return nil
	}
	
	if can, err := migration.IsExecuted(tx); nil != err {
		return err
	} else if can {
		if migration.DriverMatch(driver) {
			content, err := ioutil.ReadFile(path)
			if nil != err {
				return err
			}

			if err := tx.Exec(string(content)).Error; nil != err {
				this.logger.Info(
					"⚡️ executed migration",
					zap.String("module", migration.Module),
					zap.String("path", path),
				)

				return err
			} else {
				return migration.Save(tx)
			}
		}
	}

	return nil
}
