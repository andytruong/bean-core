package access

import (
	"context"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/access/model/dto"
	"bean/pkg/util"
)

func NewAccessModule() *AccessModule {
	return &AccessModule{
	}
}

type AccessModule struct {
	logger *zap.Logger
	db     *gorm.DB
	id     *util.Identifier
}

func (this AccessModule) Migrate(tx *gorm.DB, driver string) error {
	return nil
}

func (this *AccessModule) SessionCreate(ctx context.Context, input *dto.LoginInput) (*dto.LoginOutcome, error) {

	panic("not implemented")
}

func (this *AccessModule) SessionDelete(ctx context.Context, input *dto.LoginInput) (*dto.LogoutOutcome, error) {
	panic("not implemented")
}

func (this AccessModule) LoadSession(ctx context.Context, input *dto.ValidationInput) (*dto.ValidationOutcome, error) {
	panic("wip")
}
