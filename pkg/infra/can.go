package infra

import (
	"sync"
	"time"
	
	"go.uber.org/zap"
	
	"bean/components/conf"
	"bean/components/unique"
	"bean/pkg/access"
	"bean/pkg/integration/mailer"
	"bean/pkg/integration/s3"
	"bean/pkg/space"
)

func NewCan(path string) (*Can, error) {
	var err error
	
	this := &Can{
		mu:      &sync.Mutex{},
		bundles: beans{},
		graph:   &graph{mu: &sync.Mutex{}},
		dbs: databases{
			connections: &sync.Map{},
		},
	}
	
	// parse configuration from YAML configuration file & env variables.
	if err := conf.ParseFile(path, &this); nil != err {
		return nil, err
	}
	
	this.bundles.can = this
	this.graph.can = this
	this.dbs.config = this.Databases
	
	// setup logger
	if "dev" == this.Env {
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
	Can struct {
		Version    string                    `yaml:"version"`
		Env        string                    `yaml:"env"`
		Databases  map[string]DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig          `yaml:"http-server"`
		Bundles    BundlesConfig             `json:"bundles"`
		
		mu      *sync.Mutex
		id      *unique.Identifier
		graph   *graph
		dbs     databases
		bundles beans
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
			S3     *s3.Configuration           `yaml:"s3"`
			Mailer *mailer.MailerConfiguration `yaml:"mailer"`
		} `yaml:"integration"`
	}
)

func (this *Can) Logger() *zap.Logger {
	return this.logger
}

func (this *Can) Identifier() *unique.Identifier {
	if this.id == nil {
		this.mu.Lock()
		this.id = &unique.Identifier{}
		this.mu.Unlock()
	}
	
	return this.id
}
