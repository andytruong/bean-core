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

	this := &Container{
		mutex:   &sync.Mutex{},
		bundles: BundleList{},
	}

	// parse configuration from YAML configuration file & env variables.
	if err := conf.ParseFile(path, &this); nil != err {
		return nil, err
	}

	this.bundles.container = this
	this.dbs = connect.NewWrapper(this.Databases)

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
		Version    string                            `yaml:"version"`
		Env        string                            `yaml:"env"`
		Databases  map[string]connect.DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig                  `yaml:"http-server"`
		Bundles    BundlesConfig                     `json:"BundleList"`

		mutex   *sync.Mutex
		id      *scalar.Identifier
		dbs     *connect.Wrapper
		bundles BundleList
		logger  *zap.Logger
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

func (c *Container) Logger() *zap.Logger {
	return c.logger
}

func (c *Container) Identifier() *scalar.Identifier {
	if c.id == nil {
		c.mutex.Lock()
		c.id = &scalar.Identifier{}
		c.mutex.Unlock()
	}

	return c.id
}

func (c *Container) BundleList() BundleList {
	return c.bundles
}

func (c *Container) Bundle(i int) module.Bundle {
	list := c.BundleList()
	bundles := list.Get()

	return bundles[i]
}
