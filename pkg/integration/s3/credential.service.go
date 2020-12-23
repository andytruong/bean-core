package s3

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"

	"bean/components/util"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

type credentialService struct {
	bundle    *S3Bundle
	transport http.RoundTripper
}

func (service *credentialService) loadByApplicationId(ctx context.Context, appId string) (*model.Credentials, error) {
	cred := &model.Credentials{}

	err := service.bundle.db.WithContext(ctx).
		Where("application_id = ?", appId).
		First(&cred).
		Error

	if nil != err {
		return nil, err
	}

	return cred, nil
}

func (service *credentialService) onAppCreate(tx *gorm.DB, app *model.Application, in dto.S3ApplicationCredentialsCreateInput) error {
	cre := model.Credentials{
		ID:            service.bundle.id.MustULID(),
		ApplicationId: app.ID,
		Endpoint:      in.Endpoint,
		Bucket:        in.Bucket,
		AccessKey:     in.AccessKey,
		SecretKey:     service.encrypt(in.SecretKey),
		IsSecure:      in.IsSecure,
	}

	return tx.Create(&cre).Error
}

func (service *credentialService) onAppUpdate(tx *gorm.DB, app *model.Application, in *dto.S3ApplicationCredentialsUpdateInput) error {
	if nil == in {
		return nil
	}

	// load
	cre := &model.Credentials{}
	err := tx.Where("application_id = ?", app.ID).First(&cre).Error

	if nil == err {
		// if found -> update
		changed := false

		if nil != in.Endpoint {
			changed = true
			cre.Endpoint = *in.Endpoint
		}

		if nil != in.Bucket {
			changed = true
			cre.Bucket = *in.Bucket
		}

		if nil != in.IsSecure {
			changed = true
			cre.IsSecure = *in.IsSecure
		}

		if nil != in.AccessKey || nil != in.SecretKey {
			if nil != in.AccessKey {
				changed = true
				cre.AccessKey = *in.AccessKey
			}

			if nil != in.SecretKey {
				changed = true
				cre.SecretKey = service.encrypt(*in.SecretKey)
			}
		}

		if changed {
			return tx.Save(&cre).Error
		}

		return util.ErrorUselessInput
	} else {
		if gorm.ErrRecordNotFound != err {
			return err
		}

		if nil != in.Endpoint && nil != in.AccessKey && nil != in.SecretKey {
			// if not found -> create
			cre = &model.Credentials{
				ID:            service.bundle.id.MustULID(),
				ApplicationId: app.ID,
				Endpoint:      *in.Endpoint,
				AccessKey:     *in.AccessKey,
				SecretKey:     service.encrypt(*in.SecretKey),
				IsSecure:      false,
			}

			err := tx.Create(&cre).Error
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (service credentialService) encrypt(text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(service.bundle.config.Key))
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

func (service credentialService) decrypt(cryptoText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher([]byte(service.bundle.config.Key))
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

func (service *credentialService) client(creds *model.Credentials) (*minio.Client, error) {
	endpoint := string(creds.Endpoint)
	endpoint = strings.Replace(endpoint, "http://", "", 1)
	endpoint = strings.Replace(endpoint, "https://", "", 1)

	return minio.New(
		endpoint,
		&minio.Options{
			Creds:     credentials.NewStaticV4(creds.AccessKey, service.decrypt(creds.SecretKey), ""),
			Secure:    creds.IsSecure,
			Transport: service.transport,
		},
	)
}
