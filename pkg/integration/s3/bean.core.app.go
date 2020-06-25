package s3

import (
	"context"
	"time"

	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util/connect"
)

type CoreApplication struct {
	bean *S3IntegrationBean
}

func (this *CoreApplication) Create(ctx context.Context, in dto.S3ApplicationCreateInput) (*model.Application, error) {
	var app *model.Application

	err := connect.Transaction(
		ctx,
		this.bean.db,
		func(tx *gorm.DB) error {
			// create credentials
			cre := model.Credentials{
				ID:               this.bean.id.MustULID(),
				Version:          this.bean.id.MustULID(),
				Endpoint:         in.Credentials.Endpoint,
				EncryptedKeyPair: in.Credentials.AccessKey + " " + in.Credentials.SecretKey,
				IsSecure:         in.Credentials.IsSecure,
			}

			app = &model.Application{
				Slug:          in.Slug,
				ID:            this.bean.id.MustULID(),
				Version:       this.bean.id.MustULID(),
				IsActive:      in.IsActive,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				DeletedAt:     nil,
				CredentialsId: cre.ID,
			}

			err := this.bean.db.WithContext(ctx).Create(&app).Error
			if nil != err {
				return err
			}

			return nil
		},
	)

	if nil != err {
		return nil, err
	}

	return app, nil
}

func (this *CoreApplication) Update() {
}

func (this *CoreApplication) Delete() {
}
