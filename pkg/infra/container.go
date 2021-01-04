package infra

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"bean/components/conf"
	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
	"bean/pkg/access"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
)

func NewContainer(path string) (*Container, error) {
	var err error

	container := &Container{
		Config:     &Config{},
		mutex:      &sync.Mutex{},
		identifier: &scalar.Identifier{},
		bundles:    BundleList{},
		hook:       module.NewHook(),
	}

	// parse configuration from YAML configuration file & env variables.
	if err := conf.ParseFile(path, &container.Config); nil != err {
		return nil, err
	}

	container.bundles.container = container
	container.DBs = connect.NewWrapper(container.Config.Databases)

	// setup logger
	if container.Config.Env == "dev" {
		if container.logger, err = zap.NewDevelopment(); nil != err {
			return nil, err
		}
	} else {
		if container.logger, err = zap.NewProduction(); nil != err {
			return nil, err
		}
	}

	return container, nil
}

type (
	// Locator for most important services for system:
	// 	- Logger, ID generator, bundles, …
	//  - Database connections
	//  - HTTP server (GraphQL query interface)
	Container struct {
		Config *Config

		mutex      *sync.Mutex
		identifier *scalar.Identifier
		DBs        *connect.Wrapper
		bundles    BundleList
		logger     *zap.Logger
		hook       *module.Hook
	}

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

func (c *Container) Logger() *zap.Logger {
	return c.logger
}

func (c *Container) BundleList() BundleList {
	return c.bundles
}
