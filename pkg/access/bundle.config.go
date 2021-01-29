package access

import (
	"sync"
	"time"

	"bean/components/scalar"
)

type (
	Config struct {
		SessionTimeout time.Duration `yaml:"timeout"`
		Jwt            JwtConfig     `yaml:"jwt"`

		mutex      *sync.Mutex
		privateKey interface{}
		publicKey  interface{}
	}

	JwtConfig struct {
		Algorithm  string          `yaml:"algorithm"`
		PrivateKey scalar.FilePath `yaml:"privateKey"`
		PublicKey  scalar.FilePath `yaml:"publicKey"`
		Timeout    time.Duration   `yaml:"timeout"`
	}
)

func (cnf *Config) init() *Config {
	if nil == cnf.mutex {
		cnf.mutex = &sync.Mutex{}
	}

	if time.Duration(0) == cnf.SessionTimeout {
		cnf.SessionTimeout, _ = time.ParseDuration("128h")
	}

	// go time to validate configuration
	// …

	return cnf
}
