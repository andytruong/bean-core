package infra

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"bean/pkg/access"
	"bean/pkg/namespace"
	"bean/pkg/util"
)

func NewCan(path string) (*Can, error) {
	var err error

	this := &Can{
		mu:    &sync.Mutex{},
		beans: beans{},
		graph: &graph{mu: &sync.Mutex{}},
		dbs: databases{
			connections: &sync.Map{},
		},
	}

	// parse configuration from YAML configuration file & env variables.
	if err := util.ParseFile(path, &this); nil != err {
		return nil, err
	}

	this.beans.can = this
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
		Beans      BeansConfig               `json:"beans"`

		mu     *sync.Mutex
		id     *util.Identifier
		graph  *graph
		dbs    databases
		beans  beans
		logger *zap.Logger
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
			Playround PlayroundConfig
		} `yaml:"graphql"`
	}

	PlayroundConfig struct {
		Title   string `yaml:"title"`
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	}

	BeansConfig struct {
		Access    *access.Genetic    `yaml:"access"`
		Namespace *namespace.Genetic `yaml:"namespace"`
	}
)

func (this *Can) Logger() *zap.Logger {
	return this.logger
}

func (this *Can) Identifier() *util.Identifier {
	if this.id == nil {
		this.mu.Lock()
		this.id = &util.Identifier{}
		this.mu.Unlock()
	}

	return this.id
}
