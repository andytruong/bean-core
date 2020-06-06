package infra

import (
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"bean/pkg/access"
	"bean/pkg/infra/gql"
	"bean/pkg/util"
)

func NewContainer(path string) (*Container, error) {
	var err error

	this := &Container{}

	// parse configuration from YAML configuration file & env variables.
	if err := this.parseFile(path); nil != err {
		return nil, err
	}

	this.mu = &sync.Mutex{}

	this.modules = modules{
		container: this,
		user:      nil,
		access:    nil,
	}

	this.graph = &graph{
		container: this,
		mu:        &sync.Mutex{},
	}

	// setup logger
	if this.logger, err = zap.NewProduction(); nil != err {
		return nil, err
	}

	this.dbs = databases{
		config:      this.Databases,
		connections: &sync.Map{},
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
		Databases  map[string]DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig          `yaml:"http-server"`
		Modules    ModulesConfig             `json:"modules"`

		mu      *sync.Mutex
		id      *util.Identifier
		graph   *graph
		dbs     databases
		modules modules
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
			Playround PlayroundConfig
		} `yaml:"graphql"`
	}

	PlayroundConfig struct {
		Title   string `yaml:"title"`
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	}

	ModulesConfig struct {
		Access *access.Config `yaml:"access"`
	}
)

func (this *Container) Logger() *zap.Logger {
	return this.logger
}

func (this *Container) parseFile(path string) error {
	content, err := util.ParseFile(path)
	if nil != err {
		return err
	}

	return yaml.Unmarshal(content, &this)
}

func (this *Container) ListenAndServe() error {
	router := mux.NewRouter()

	cnf := gql.Config{Resolvers: this.graph}
	schema := gql.NewExecutableSchema(cnf)
	hdl := handler.NewDefaultServer(schema)
	router.HandleFunc("/query", hdl.ServeHTTP)

	if this.HttpServer.GraphQL.Playround.Enabled {
		hdl := playground.Handler(this.HttpServer.GraphQL.Playround.Title, "/query")
		router.Handle(this.HttpServer.GraphQL.Playround.Path, hdl)
	}

	server := http.Server{
		Addr:              this.HttpServer.Address,
		Handler:           router,
		ReadTimeout:       this.HttpServer.ReadTimeout,
		ReadHeaderTimeout: this.HttpServer.ReadTimeout,
		WriteTimeout:      this.HttpServer.WriteTimeout,
		IdleTimeout:       this.HttpServer.ReadTimeout,
	}

	this.logger.Info("🚀 HTTP server is running", zap.String("address", this.HttpServer.Address))

	return server.ListenAndServe()
}

func (this *Container) Identifier() *util.Identifier {
	if this.id == nil {
		this.mu.Lock()
		this.id = &util.Identifier{}
		this.mu.Unlock()
	}

	return this.id
}
