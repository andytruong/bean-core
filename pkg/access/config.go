package access

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"

	"bean/pkg/util"
)

type (
	Config struct {
		SessionTimeout time.Duration `yaml:"sessionTimeout"`
		Jwt            JwtConfig     `yaml:"jwt"`

		mutex      *sync.Mutex
		privateKey interface{}
	}

	JwtConfig struct {
		Algorithm  string `yaml:"algorithm"`
		PrivateKey string `yaml:"privateKey"`
		PublicKey  string `yaml:"publicKey"`
		Timeout    time.Duration
	}
)

func (this *Config) init() *Config {
	this.mutex = &sync.Mutex{}

	// go time to validate configuration
	// …

	return this
}

func (this *Config) signMethod() jwt.SigningMethod {
	switch this.Jwt.Algorithm {
	case "RS512":
		return jwt.SigningMethodRS512

	default:
		panic(util.ErrorToBeImplemented)
	}
}

func (this *Config) signKey() (interface{}, error) {
	if nil == this.privateKey {
		this.mutex.Lock()
		defer this.mutex.Unlock()

		file, err := ioutil.ReadFile(this.Jwt.PrivateKey)
		if nil != err {
			return nil, err
		}

		block, _ := pem.Decode(file)
		this.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)

		if err != nil {
			return nil, err
		}
	}

	return this.privateKey, nil
}
