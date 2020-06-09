package config

import (
	"context"
	"path"
	"runtime"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"bean/pkg/config/model"
	"bean/pkg/config/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
	"bean/pkg/util/migrate"
)

func NewConfigBean(id *util.Identifier, logger *zap.Logger) *ConfigBean {
	return &ConfigBean{
		id:     nil,
		logger: nil,
	}
}

type ConfigBean struct {
	id     *util.Identifier
	logger *zap.Logger
}

func (this ConfigBean) Migrate(tx *gorm.DB, driver string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil
	}

	runner := migrate.Runner{
		Tx:     tx,
		Logger: this.logger,
		Driver: driver,
		Bean:   "config",
		Dir:    path.Dir(filename) + "/model/migration/",
	}

	return runner.Run()
}

func (this ConfigBean) Dependencies() []util.Bean {
	panic("implement me")
}

func (this ConfigBean) BucketCreate(ctx context.Context, tx *gorm.DB, input dto.BucketCreateInput) (*dto.BucketMutationOutcome, error) {
	bucket := &model.ConfigBucket{
		Id:          this.id.MustULID(),
		Version:     this.id.MustULID(),
		Slug:        util.NotNilString(input.Slug, this.id.MustULID()),
		Title:       *input.Title,
		Description: input.Description,
		Access:      "777",
		HostId:      input.HostId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if nil != input.Access {
		bucket.Access = *input.Access
	}

	err := tx.Table(connect.TableNamespaceConfig).Save(&bucket).Error
	if nil != err {
		return nil, err
	}

	return &dto.BucketMutationOutcome{
		Errors: nil,
		Bucket: bucket,
	}, nil
}

func (this ConfigBean) BucketUpdate(ctx context.Context, tx *gorm.DB, input dto.BucketUpdateInput) (*dto.BucketMutationOutcome, error) {
	bucket, err := this.BucketLoad(ctx, tx, input.Id)
	if nil != err {
		return nil, err
	}

	if bucket.Version != input.Version {
		return nil, util.ErrorVersionConflict
	}

	changed := false
	if input.Title != nil {
		if bucket.Title != *input.Title {
			changed = true
			bucket.Title = *input.Title
		}
	}

	if input.Description != nil {
		if bucket.Description != input.Description {
			changed = true
			bucket.Description = input.Description
		}
	}

	if input.Access != nil {
		if bucket.Access != *input.Access {
			changed = true
			bucket.Access = *input.Access
		}
	}

	if changed {
		bucket.Version = this.id.MustULID()
		err = tx.
			Table(connect.TableNamespaceConfig).
			Save(&bucket).
			Error
		if nil != err {
			return nil, err
		}
	}

	return &dto.BucketMutationOutcome{
		Errors: nil,
		Bucket: bucket,
	}, nil
}

func (this ConfigBean) BucketLoad(ctx context.Context, db *gorm.DB, id string) (*model.ConfigBucket, error) {
	bucket := &model.ConfigBucket{}

	err := db.
		Table(connect.TableConfigBucket).
		First(&bucket, "id = ?", id).
		Error
	if nil != err {
		return nil, err
	}

	return bucket, nil
}
