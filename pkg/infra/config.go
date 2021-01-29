package infra

import (
	"time"

	"bean/components/connect"
	"bean/pkg/access"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
)

type (
	Config struct {
		Version    string                            `yaml:"version"`
		Env        string                            `yaml:"env"`
		Databases  map[string]connect.DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig                  `yaml:"http-server"`
		Bundles    BundlesConfig                     `yaml:"bundles"`
	}

	HttpServerConfig struct {
		Address      string        `yaml:"address"`
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
		IdleTimeout  time.Duration `yaml:"idleTimeout"`
		GraphQL      struct {
			Introspection bool `yaml:"introspection"`
			Transports    struct {
				Post      bool `yaml:"post"`
				Websocket struct {
					KeepAlivePingInterval time.Duration `yaml:"keepAlivePingInterval"`
				} `yaml:"websocket"`
			} `yaml:"transports"`
			Playround PlayroundConfig `yaml:"playround"`
		} `yaml:"graphql"`
	}

	PlayroundConfig struct {
		Title   string `yaml:"title"`
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	}

	BundlesConfig struct {
		Access      *access.Config `yaml:"access"`
		Space       *space.Config  `yaml:"space"`
		Integration struct {
			S3     *s3.Config     `yaml:"s3"`
			Mailer *mailer.Config `yaml:"mailer"`
		} `yaml:"integration"`
	}
)
