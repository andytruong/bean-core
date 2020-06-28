package s3

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"

	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
	"bean/pkg/util"
	"bean/pkg/util/connect"
)

type CoreApplication struct {
	bean *S3IntegrationBean
}

func (this *CoreApplication) Load(ctx context.Context, id string) (*model.Application, error) {
	app := &model.Application{}

	err := this.bean.db.
		WithContext(ctx).
		Table(connect.TableIntegrationS3).
		Where("id = ?", id).
		First(&app).
		Error
	if nil != err {
		return nil, err
	}

	return app, nil
}

func (this *CoreApplication) Create(ctx context.Context, in dto.S3ApplicationCreateInput) (*dto.S3ApplicationMutationOutcome, error) {
	var app *model.Application

	err := connect.Transaction(
		ctx,
		this.bean.db,
		func(tx *gorm.DB) error {
			app = &model.Application{
				Slug:      in.Slug,
				ID:        this.bean.id.MustULID(),
				Version:   this.bean.id.MustULID(),
				IsActive:  in.IsActive,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			}

			err := tx.Table(connect.TableIntegrationS3).Create(&app).Error
			if nil != err {
				return err
			} else {
				// create credentials
				cre := model.Credentials{
					ID:               this.bean.id.MustULID(),
					Version:          this.bean.id.MustULID(),
					ApplicationId:    app.ID,
					Endpoint:         in.Credentials.Endpoint,
					EncryptedKeyPair: in.Credentials.AccessKey + " " + this.encrypt(in.Credentials.SecretKey),
					IsSecure:         in.Credentials.IsSecure,
				}

				err := tx.Table(connect.TableIntegrationS3Credentials).Create(&cre).Error
				if nil != err {
					return err
				}
			}

			return nil
		},
	)

	if nil != err {
		return nil, err
	}

	return &dto.S3ApplicationMutationOutcome{App: app, Errors: nil}, nil
}

func (this *CoreApplication) Update(ctx context.Context, in dto.S3ApplicationUpdateInput) (*dto.S3ApplicationMutationOutcome, error) {
	app, err := this.Load(ctx, in.Id)

	if nil != err {
		return nil, err
	} else if app.Version != in.Version {
		return nil, util.ErrorVersionConflict
	}

	changed := false
	if nil != in.IsActive {
		if app.IsActive != *in.IsActive {
			app.IsActive = *in.IsActive
			changed = true
		}
	}

	if nil != in.Slug {
		if app.Slug != *in.Slug {
			app.Slug = *in.Slug
			changed = true
		}
	}

	if deletedAt, ok := ctx.Value("bean.integration-s3.delete").(time.Time); ok {
		app.DeletedAt = &deletedAt
		changed = true
	}

	if !changed {
		if nil == in.Credentials {
			return nil, util.ErrorUselessInput
		}
	}

	app.Version = this.bean.id.MustULID()
	app.UpdatedAt = time.Now()
	err = connect.Transaction(
		ctx,
		this.bean.db,
		func(tx *gorm.DB) error {
			err := tx.Save(&app).Error
			if nil != err {
				return err
			}

			if nil != in.Credentials {
				var cred *model.Credentials

				return tx.
					Table(connect.TableIntegrationS3Credentials).
					First(cred, "application_id = ?", app.ID).
					Error
			}

			return nil
		},
	)

	return &dto.S3ApplicationMutationOutcome{App: app, Errors: nil}, err
}

func (this *CoreApplication) Delete(ctx context.Context, in dto.S3ApplicationDeleteInput) (*dto.S3ApplicationMutationOutcome, error) {
	ctx = context.WithValue(ctx, "bean.integration-s3.delete", time.Now())

	return this.Update(ctx, dto.S3ApplicationUpdateInput{
		Id:       in.Id,
		Version:  in.Version,
		IsActive: util.NilBool(true),
	})
}

func (this CoreApplication) encrypt(text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(this.bean.genetic.Key)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(cipherText)
}

func (this CoreApplication) decrypt(cryptoText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher(this.bean.genetic.Key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return fmt.Sprintf("%s", cipherText)
}
