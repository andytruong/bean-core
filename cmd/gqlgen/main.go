package main

import (
	"fmt"
	"os"
	
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}
	
	err = api.Generate(cfg, api.AddPlugin(MyPlugin{}))
	
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

type MyPlugin struct {
}

func (this MyPlugin) Name() string {
	return "bean"
}

func (this MyPlugin) GenerateCode(cfg *codegen.Data) error {
	fmt.Println("GenerateCode", cfg)
	
	panic("implement me")
}
