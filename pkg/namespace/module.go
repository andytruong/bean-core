package namespace

import (
	"path"
	"runtime"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/util"
	"bean/pkg/util/migrate"
)

func NewNamespaceModule(db *gorm.DB, logger *zap.Logger, id *util.Identifier) (*NamespaceModule, error) {
	module := &NamespaceModule{
		logger: logger,
		Mutation: NamespaceMutationResolver{
			db: db,
			id: id,
		},
		Query: NamespaceQueryResolver{
			db: db,
		},
	}

	return module, nil
}

type NamespaceModule struct {
	logger   *zap.Logger
	Mutation NamespaceMutationResolver
	Query    NamespaceQueryResolver
}

func (this NamespaceModule) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Module: "namespace",
		Dir:    path.Dir(filename) + "/migration/",
	}

	return runner.Run()
}
