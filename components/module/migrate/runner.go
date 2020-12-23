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
	Bundle string
	Dir    string
}

func (runner Runner) Run() error {
	return filepath.Walk(
		runner.Dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.HasSuffix(path, ".sql") {
					return runner.installFile(path)
				}
			}

			return nil
		},
	)
}

func (runner Runner) installFile(file string) error {
	migration := NewMigration(runner.Bundle, file)
	path := migration.RealPath()

	if !migration.DriverMatch(runner.Driver) {
		runner.Logger.Debug(
			"👉 driver unmatched",
			zap.String("bean", migration.Bundle),
			zap.String("path", path),
		)

		return nil
	}

	if can, err := migration.IsExecuted(runner.Tx); nil != err {
		return err
	} else if can {
		if migration.DriverMatch(runner.Driver) {
			content, err := ioutil.ReadFile(path)
			if nil != err {
				return err
			}

			if err := runner.Tx.Exec(string(content)).Error; nil != err {
				runner.Logger.Info(
					"⚡️ executed migration",
					zap.String("bean", migration.Bundle),
					zap.String("path", path),
				)

				return err
			} else {
				return migration.Save(runner.Tx)
			}
		}
	}

	return nil
}
