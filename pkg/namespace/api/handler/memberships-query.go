package handler

import (
	"context"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/namespace/model"
	"bean/pkg/namespace/model/dto"
	"bean/pkg/util/api"
)

type MembershipsQueryHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (this MembershipsQueryHandler) Memberships(
	ctx context.Context,
	first int, after *string, filters dto.MembershipsFilter, sort *api.Sorts,
) (*model.MembershipConnection, error) {
	panic("implement me")
}
