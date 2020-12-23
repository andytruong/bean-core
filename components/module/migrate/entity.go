package migrate

import (
	"path"
	"runtime"
	"strings"
	"time"

	"gorm.io/gorm"
)

func NewMigration(bundleName string, name string) Migration {
	this := Migration{
		Bundle:    bundleName,
		Name:      name,
		CreatedAt: time.Now(),
	}

	this.Name = strings.TrimPrefix(this.Name, RootDirectory())
	this.Name = strings.TrimPrefix(this.Name, "/")

	return this
}

type Migration struct {
	Bundle    string `gorm:"unique_index:bundle_unique_schema"`
	Name      string `gorm:"unique_index:bundle_unique_schema"`
	CreatedAt time.Time
}

func (migration Migration) RealPath() string {
	return RootDirectory() + "/" + migration.Name
}

func (migration Migration) DriverMatch(driver string) bool {
	return strings.HasSuffix(migration.Name, "."+driver+".sql")
}

func (migration *Migration) IsExecuted(tx *gorm.DB) (bool, error) {
	var count int64

	err := tx.
		Model(&Migration{}).
		Where(&Migration{Bundle: migration.Bundle, Name: migration.Name}).
		Count(&count).
		Error

	if nil != err {
		return false, err
	}

	return count == 0, nil
}

func (migration *Migration) Save(tx *gorm.DB) error {
	return tx.Create(migration).Error
}

func RootDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	dir = path.Dir(dir)
	dir = path.Dir(dir)
	dir = path.Dir(dir)

	return dir
}
