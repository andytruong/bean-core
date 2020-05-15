package namespace

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/user/service"
	"bean/pkg/util"
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

func (this NamespaceModule) Install(tx *gorm.DB, driver string) error {
	api := service.NewUserInstallAPI(this.logger)

	return api.Run(tx, driver)
}
