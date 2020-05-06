package infra

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"bean/pkg/infra/gql"
	"bean/pkg/util"
)

func NewContainer(path string) (*Container, error) {
	var err error

	this := &Container{}
	this.mu = &sync.Mutex{}

	this.modules = modules{
		container: this,
		user:      nil,
		access:    nil,
	}
	
	this.gql = resolvers{
		container: this,
		mu:        &sync.Mutex{},
	}

	// setup logger
	if this.Logger, err = zap.NewProduction(); nil != err {
		return nil, err
	}

	if err := this.parseFile(path); nil != err {
		return nil, err
	}

	return this, nil
}

type (
	// Locator for most important services for system:
	// 	- Logger
	//  - Database connections
	//  - HTTP server (GraphQL query interface)
	Container struct {
		Version    string `yaml:"version"`
		Logger     *zap.Logger
		DB         *gorm.DB
		Databases  map[string]DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig          `yaml:"http-server"`

		mu      *sync.Mutex
		id      *util.Identifier
		gql     resolvers
		modules modules
	}

	DatabaseConfig struct {
		Driver string `yaml:""`
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
)

func (this *Container) parseFile(path string) error {
	raw, err := ioutil.ReadFile(path)
	if nil != err {
		return err
	} else if content, err := this.replaceEnvVariables(raw); nil != err {
		return err
	} else if err := yaml.Unmarshal(content, &this); nil != err {
		return err
	}

	return nil
}

func (this *Container) replaceEnvVariables(inBytes []byte) ([]byte, error) {
	if envRegex, err := regexp.Compile(`\${[0-9A-Za-z_]+(:((\${[^}]+})|[^}])+)?}`); err != nil {
		return nil, err
	} else if escapedEnvRegex, err := regexp.Compile(`\${({[0-9A-Za-z_]+(:((\${[^}]+})|[^}])+)?})}`); err != nil {
		return nil, err
	} else {
		replaced := envRegex.ReplaceAllFunc(inBytes, func(content []byte) []byte {
			var value string
			if len(content) > 3 {
				if colonIndex := bytes.IndexByte(content, ':'); colonIndex == -1 {
					value = os.Getenv(string(content[2 : len(content)-1]))
				} else {
					targetVar := content[2:colonIndex]
					defaultVal := content[colonIndex+1 : len(content)-1]

					value = os.Getenv(string(targetVar))
					if len(value) == 0 {
						value = string(defaultVal)
					}
				}
			}
			return []byte(value)
		})

		return escapedEnvRegex.ReplaceAll(replaced, []byte("$$$1")), nil
	}
}

func (this *Container) ListenAndServe() error {
	router := mux.NewRouter()
	router.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		config := gql.Config{Resolvers: this.gql.getRoot()}
		schema := gql.NewExecutableSchema(config)
		server := handler.NewDefaultServer(schema)
		server.ServeHTTP(w, r)
	})

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

	this.Logger.Info("🚀 HTTP server is running", zap.String("address", this.HttpServer.Address))

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
