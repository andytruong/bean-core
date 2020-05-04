package infra

import (
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"bean/pkg/infra/gql"
)

func NewContainer(path string) (*Container, error) {
	raw, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}

	content := os.ExpandEnv(string(raw))
	this := &Container{
		mu: &sync.Mutex{},
	}

	if err := yaml.Unmarshal([]byte(content), &this); nil != err {
		return nil, err
	}

	// setup gql-resolvers
	if false {
		this.gql = struct {
			root     *rootResolver
			query    *queryResolver
			mutation *mutationResolver
			session  *sessionResolver
		}{
			root:  &rootResolver{container: this},
			query: &queryResolver{AccessQueryResolver: nil},
			mutation: &mutationResolver{
				AccessMutationResolver: nil,
				UserMutationResolver:   nil,
			},
			session: &sessionResolver{
				container: this,
			},
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
		Version    string `yaml:"version"`
		Logger     zap.Logger
		DB         *gorm.DB
		Databases  map[string]DatabaseConfig `yaml:"databases"`
		HttpServer HttpServerConfig          `yaml:"http-server"`

		mu  *sync.Mutex
		gql struct {
			root     *rootResolver
			query    *queryResolver
			mutation *mutationResolver
			session  *sessionResolver
		}
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

func (this *Container) ListenAndServe() error {
	router := mux.NewRouter()
	router.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		config := &gql.Config{Resolvers: this.gql.root}
		schema := gql.NewExecutableSchema(*config)
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
