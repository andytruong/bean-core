package util

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type MigrationRunner struct {
	Tx     *gorm.DB
	Logger *zap.Logger
	Driver string
	Module string
	Dir    string
}

func (this MigrationRunner) Run() error {
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

func (this MigrationRunner) installFile(file string) error {
	migration := NewMigration(this.Module, file)
	path := migration.RealPath()

	if !migration.DriverMatch(this.Driver) {
		this.Logger.Info(
			"👉 driver unmatched",
			zap.String("module", migration.Module),
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
					zap.String("module", migration.Module),
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

func NewMigration(module string, name string) Migration {
	this := Migration{
		Module:    module,
		Name:      name,
		CreatedAt: time.Now(),
	}

	this.Name = strings.TrimPrefix(this.Name, RootDirectory())
	this.Name = strings.TrimPrefix(this.Name, "/")

	return this
}

type (
	Module interface {
		Install(tx *gorm.DB, driver string) error
	}

	Migration struct {
		Module    string `gorm:"unique_index:module_unique_schema"`
		Name      string `gorm:"unique_index:module_unique_schema"`
		CreatedAt time.Time
	}
)

func RootDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	pkg := path.Dir(dir)
	root := path.Dir(pkg)

	return root
}

func (this Migration) RealPath() string {
	return RootDirectory() + "/" + this.Name
}

func (this Migration) DriverMatch(driver string) bool {
	return strings.HasSuffix(this.Name, "."+driver+".sql")
}

func (this *Migration) IsExecuted(tx *gorm.DB) (bool, error) {
	count := 0

	err := tx.
		Model(&Migration{}).
		Where(&Migration{Module: this.Module, Name: this.Name}).
		Count(&count).
		Error

	if nil != err {
		return false, err
	}

	return count == 0, nil
}

func (this *Migration) Save(tx *gorm.DB) error {
	return tx.Create(this).Error
}
