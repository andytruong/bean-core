package access

import (
	"time"

	"bean/components/scalar"
)

type (
	Config struct {
		SessionTimeout time.Duration `yaml:"timeout"`
		Jwt            JwtConfig     `yaml:"jwt"`
	}

	JwtConfig struct {
		Algorithm  string          `yaml:"algorithm"`
		PrivateKey scalar.FilePath `yaml:"privateKey"`
		PublicKey  scalar.FilePath `yaml:"publicKey"`
		Timeout    time.Duration   `yaml:"timeout"`
	}
)

func (cnf *Config) init() *Config {
	if time.Duration(0) == cnf.SessionTimeout {
		cnf.SessionTimeout, _ = time.ParseDuration("128h")
	}

	// go time to validate configuration
	// …

	return cnf
}
