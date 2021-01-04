package connect

import (
	"context"
	"database/sql"
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

func NewMigration(bundleName string, name string) Migration {
	this := Migration{
		Bundle:    bundleName,
		Name:      name,
		CreatedAt: time.Now(),
	}

	this.Name = strings.TrimPrefix(this.Name, scalar.RootDirectory())
	this.Name = strings.TrimPrefix(this.Name, "/")

	return this
}

type (
	Migration struct {
		Bundle    string `gorm:"unique_index:bundle_unique_schema"`
		Name      string `gorm:"unique_index:bundle_unique_schema"`
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
	tx := db.WithContext(ctx).Begin()
	ctx = DBToContext(ctx, db)

	// create migration table if not existing
	if !tx.Migrator().HasTable(Migration{}) {
		if err := tx.Migrator().CreateTable(Migration{}); nil != err {
			tx.Rollback()
			return err
		}
	}

	for _, bundle := range bundles {
		if err := bundle.Migrate(ctx, driver); nil != err {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (mig Migration) realPath() string {
	return scalar.RootDirectory() + "/" + mig.Name
}

func (mig Migration) driverMatch(driver string) bool {
	return strings.HasSuffix(mig.Name, "."+driver+".sql")
}

func (mig *Migration) isExecuted(db *gorm.DB) (bool, error) {
	var count int64

	err := db.
		Model(&Migration{}).
		Where(&Migration{Bundle: mig.Bundle, Name: mig.Name}).
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
