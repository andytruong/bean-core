package migrate

import (
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

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

type Migration struct {
	Module    string `gorm:"unique_index:module_unique_schema"`
	Name      string `gorm:"unique_index:module_unique_schema"`
	CreatedAt time.Time
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

func RootDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	dir = path.Dir(dir)
	dir = path.Dir(dir)
	dir = path.Dir(dir)

	return dir
}
