package infra

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"bean/components/conf"
	"bean/components/module"
	"bean/components/unique"
	"bean/pkg/access"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
	"bean/pkg/user"
)

func NewContainer(path string) (*Container, error) {
	var err error

	this := &Container{
		mutex:   &sync.Mutex{},
		bundles: bundles{},
		dbs: databases{
			connections: &sync.Map{},
		},
	}

	// parse configuration from YAML configuration file & env variables.
	if err := conf.ParseFile(path, &this); nil != err {
		return nil, err
	}

	this.bundles.container = this
	this.dbs.config = this.Databases

	// setup logger
	if this.Env == "dev" {
		if this.logger, err = zap.NewDevelopment(); nil != err {
			return nil, err
		}
	} else {
		if this.logger, err = zap.NewProduction(); nil != err {
			return nil, err
		}
	}

	return this, nil
}

type (
	// Locator for most important services for system:
	// 	- Logger
	//  - Database connections
	//  - HTTP server (GraphQL query interface)
	Container struct {
		Version    string                    `yaml:"version"`
		Env        string                    `yaml:"env"`
		Databases  map[string]DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig          `yaml:"http-server"`
		Bundles    BundlesConfig             `json:"bundles"`

		mutex   *sync.Mutex
		id      *unique.Identifier
		dbs     databases
		bundles bundles
		logger  *zap.Logger
	}

	DatabaseConfig struct {
		Driver string `yaml:"driver"`
		Url    string `yaml:"url"`
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
		Access      *access.AccessConfiguration `yaml:"access"`
		Space       *space.SpaceConfiguration   `yaml:"space"`
		Integration struct {
			S3     *s3.S3Configuration         `yaml:"s3"`
			Mailer *mailer.MailerConfiguration `yaml:"mailer"`
		} `yaml:"integration"`
	}
)

func (container *Container) Logger() *zap.Logger {
	return container.logger
}

func (container *Container) Identifier() *unique.Identifier {
	if container.id == nil {
		container.mutex.Lock()
		container.id = &unique.Identifier{}
		container.mutex.Unlock()
	}

	return container.id
}

func (container *Container) BundleList() []module.Bundle {
	return container.bundles.List()
}

func (container *Container) Bundle(i int) module.Bundle {
	bundles := container.BundleList()

	return bundles[i]
}

// TODO: Generate this code
func (container *Container) BundlePath(bundle module.Bundle) string {
	switch bundle.(type) {
	case *user.UserBundle:
		return "User"

	case *space.SpaceBundle:
		return "Space"

	case *access.AccessBundle:
		return "Access"

	case *s3.S3Bundle:
		return "S3"

	case *mailer.MailerBundle:
		return "Mailer"
	}

	panic("unknown bundle")
}
