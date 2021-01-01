package s3

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"

	"bean/components/connect"
	"bean/components/util"
	appModel "bean/pkg/app/model"
	"bean/pkg/integration/s3/model"
	"bean/pkg/integration/s3/model/dto"
)

const (
	credentialsConfigSlug   = `bean.s3.credentials.schema.v1`
	credentialsConfigSchema = `{
		"type":       "object",
		"required":   ["kind", "value"],
		"properties": {
			"kind":      { "type": "string", "enum": ["file_extensions", "rate_limit"] },
			"value":     {
				"type":     "string",
				"examples": [
					"file extension: pdf txt zip gz",
					"rate limit:     1MB/user/hour"
				]
			}
		}
	}`
)

type credentialService struct {
	bundle    *S3Bundle
	transport http.RoundTripper
}

func (srv *credentialService) loadByApplicationId(ctx context.Context, appId string) (*model.Credentials, error) {
	cred := &model.Credentials{}
	db := connect.ContextToDB(ctx)
	err := db.WithContext(ctx).Where("application_id = ?", appId).First(&cred).Error
	if nil != err {
		return nil, err
	}

	return cred, nil
}

func (srv *credentialService) onAppCreate(ctx context.Context, app *appModel.Application, in dto.S3ApplicationCredentialsCreateInput) error {
	db := connect.ContextToDB(ctx)
	cre := model.Credentials{
		ID:            srv.bundle.idr.MustULID(),
		ApplicationId: app.ID,
		Endpoint:      in.Endpoint,
		Bucket:        in.Bucket,
		AccessKey:     in.AccessKey,
		SecretKey:     srv.encrypt(in.SecretKey),
		IsSecure:      in.IsSecure,
	}

	// create app->config->credentials

	return db.Create(&cre).Error
}

func (srv *credentialService) onAppUpdate(ctx context.Context, app *appModel.Application, in *dto.S3ApplicationCredentialsUpdateInput) error {
	if nil == in {
		return nil
	}

	// load
	db := connect.ContextToDB(ctx)
	cre := &model.Credentials{}
	err := db.Where("application_id = ?", app.ID).First(&cre).Error

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
				cre.SecretKey = srv.encrypt(*in.SecretKey)
			}
		}

		if changed {
			return db.Save(&cre).Error
		}

		return util.ErrorUselessInput
	} else {
		if gorm.ErrRecordNotFound != err {
			return err
		}

		if nil != in.Endpoint && nil != in.AccessKey && nil != in.SecretKey {
			// if not found -> create
			cre = &model.Credentials{
				ID:            srv.bundle.idr.MustULID(),
				ApplicationId: app.ID,
				Endpoint:      *in.Endpoint,
				AccessKey:     *in.AccessKey,
				SecretKey:     srv.encrypt(*in.SecretKey),
				IsSecure:      false,
			}

			err := db.Create(&cre).Error
			if nil != err {
				return err
			}
		}
	}

	return nil
}

func (srv credentialService) encrypt(text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(srv.bundle.cnf.Key))
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

func (srv credentialService) decrypt(cryptoText string) string {
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)
	block, err := aes.NewCipher([]byte(srv.bundle.cnf.Key))
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

	return string(cipherText)
}

func (srv *credentialService) client(creds *model.Credentials) (*minio.Client, error) {
	endpoint := string(creds.Endpoint)
	endpoint = strings.Replace(endpoint, "http://", "", 1)
	endpoint = strings.Replace(endpoint, "https://", "", 1)

	return minio.New(
		endpoint,
		&minio.Options{
			Creds:     credentials.NewStaticV4(creds.AccessKey, srv.decrypt(creds.SecretKey), ""),
			Secure:    creds.IsSecure,
			Transport: srv.transport,
		},
	)
}
