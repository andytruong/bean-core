package infra

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"bean/components/conf"
	"bean/components/connect"
	"bean/components/module"
	"bean/components/scalar"
)

func NewContainer(path string) (*Container, error) {
	var err error

	container := &Container{
		Config:  &Config{},
		mutex:   &sync.Mutex{},
		idr:     &scalar.Identifier{},
		bundles: bundleList{},
		hook:    module.NewHook(),
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
		DBs    *connect.Wrapper

		mutex   *sync.Mutex
		idr     *scalar.Identifier
		logger  *zap.Logger
		hook    *module.Hook
		bundles bundleList
	}
)

func (c *Container) Logger() *zap.Logger {
	return c.logger
}

func (c *Container) BundleList() []module.Bundle {
	return c.bundles.Get()
}

func (c *Container) HttpServer() http.Server {
	router := mux.NewRouter()
	handler := graphqlHttpHandler{container: c}

	return http.Server{
		Addr:              c.Config.HttpServer.Address,
		Handler:           handler.Get(router),
		ReadTimeout:       c.Config.HttpServer.ReadTimeout,
		ReadHeaderTimeout: c.Config.HttpServer.ReadTimeout,
		WriteTimeout:      c.Config.HttpServer.WriteTimeout,
		IdleTimeout:       c.Config.HttpServer.ReadTimeout,
	}
}
