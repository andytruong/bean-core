package connect

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"bean/components/module"
	"bean/components/scalar"
)

func NewMigration(bundleName string, fileName string) Migration {
	this := Migration{
		Bundle:    bundleName,
		FileName:  fileName,
		CreatedAt: time.Now(),
	}

	this.FileName = strings.TrimPrefix(this.FileName, scalar.RootDirectory())
	this.FileName = strings.TrimPrefix(this.FileName, "/")

	return this
}

type (
	Migration struct {
		Bundle    string `gorm:"unique_index:bundle_unique_schema"`
		FileName  string `gorm:"unique_index:bundle_unique_schema"`
		CreatedAt time.Time
	}

	Runner struct {
		Logger *zap.Logger
		Driver string
		Bundle string
		Dir    string
	}
)

func dbDriver(db *sql.DB) string {
	switch db.Driver().(type) {
	case *sqlite3.SQLiteDriver:
		return SQLite

	default:
		return Postgres
	}
}

func Migrate(ctx context.Context, bundles []module.Bundle, db *gorm.DB) error {
	con, err := db.DB()
	if nil != err {
		return err
	}

	driver := dbDriver(con)
	fmt.Println("DRIVER: ", driver)
	tx := db.WithContext(ctx).Begin()
	ctx = DBToContext(ctx, tx)

	// create migration table if not existing
	if !tx.Migrator().HasTable(Migration{}) {
		if err := tx.Migrator().CreateTable(Migration{}); nil != err {
			tx.Rollback()
			return err
		}
	}

	for _, bundle := range bundles {
		if err := migrate(ctx, driver, bundle); nil != err {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func migrate(ctx context.Context, driver string, bundle module.Bundle) error {
	dependencies := bundle.Dependencies()
	if nil != dependencies {
		for _, dependency := range dependencies {
			err := migrate(ctx, driver, dependency)
			if nil != err {
				return err
			}
		}
	}

	return bundle.Migrate(ctx, driver)
}

func (mig Migration) realPath() string {
	return scalar.RootDirectory() + "/" + mig.FileName
}

func (mig Migration) driverMatch(driver string) bool {
	return strings.HasSuffix(mig.FileName, "."+driver+".sql")
}

func (mig *Migration) isExecuted(db *gorm.DB) (bool, error) {
	var count int64

	err := db.
		Model(&Migration{}).
		Where(&Migration{Bundle: mig.Bundle, FileName: mig.FileName}).
		Count(&count).
		Error

	if nil != err {
		return false, err
	}

	return count == 0, nil
}

// Context must contain Database transaction
func (runner Runner) Run(ctx context.Context) error {
	return filepath.Walk(
		runner.Dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				if strings.HasSuffix(path, ".sql") {
					return runner.execute(ctx, path)
				}
			}

			return nil
		},
	)
}

func (runner Runner) execute(ctx context.Context, file string) error {
	db := ContextToDB(ctx)
	migration := NewMigration(runner.Bundle, file)
	path := migration.realPath()

	if !migration.driverMatch(runner.Driver) {
		runner.Logger.Debug(
			"👉 driver unmatched",
			zap.String("bean", migration.Bundle),
			zap.String("path", path),
		)

		return nil
	}

	if can, err := migration.isExecuted(db); nil != err {
		return err
	} else if can {
		if migration.driverMatch(runner.Driver) {
			content, err := ioutil.ReadFile(path)
			if nil != err {
				return err
			}

			if err := db.Exec(string(content)).Error; nil != err {
				runner.Logger.Info(
					"⚡️ executed migration",
					zap.String("bean", migration.Bundle),
					zap.String("path", path),
				)

				return err
			} else {
				return db.Create(migration).Error
			}
		}
	}

	return nil
}
