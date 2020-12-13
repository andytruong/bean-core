package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	. "github.com/99designs/gqlgen/plugin/resolvergen"
	"github.com/pkg/errors"
	
	"bean/pkg/infra"
	"bean/pkg/util"
)

func main() {
	path := os.Getenv("CONFIG")
	if "" == path {
		err := errors.Wrap(util.ErrorConfig, "missing env CONFIG")
		panic(err)
	}
	
	container, err := infra.NewContainer(path)
	if nil != err {
		panic("failed creating container: " + err.Error())
	} else if cfg, err := config.LoadConfigFromDefaultLocations(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	} else {
		err = api.Generate(cfg, api.AddPlugin(MyPlugin{container: container}))
		
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(3)
		}
	}
}

type MyPlugin struct {
	container *infra.Container
}

func (this MyPlugin) Name() string {
	return "bean"
}

func (this MyPlugin) GenerateCode(data *codegen.Data) error {
	file := File{}
	
	for _, o := range data.Objects {
		if o.HasResolvers() {
			file.Objects = append(file.Objects, o)
		}
		for _, f := range o.Fields {
			if !f.IsResolver {
				continue
			}
			
			file.Resolvers = append(file.Resolvers, this.newResolver(o, f))
		}
	}
	
	resolverBuild := &ResolverBuild{
		File:         &file,
		PackageName:  data.Config.Resolver.Package,
		ResolverType: data.Config.Resolver.Type,
		HasRoot:      true,
	}
	
	options := templates.Options{
		Filename:    "pkg/infra/graphql-resolvers.go",
		PackageName: "infra",
		FileNotice:  `// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.`,
		Data:        resolverBuild,
		Packages:    data.Config.Packages,
	}
	
	return templates.Render(options)
}

func (this MyPlugin) newResolver(o *codegen.Object, f *codegen.Field) *Resolver {
	resolver := &Resolver{
		Object:         o,
		Field:          f,
		Implementation: this.resolverBody(o, f),
	}
	
	return resolver
}

func (this MyPlugin) resolverBody(o *codegen.Object, f *codegen.Field) string {
	implementation := `panic("no implementation")`
	
	for _, bundle := range this.container.BundleList() {
		resolvers := bundle.GraphqlResolver()
		if nil == resolvers {
			continue
		}
		
		if objResolver, ok := resolvers[o.Name]; ok {
			if fieldResolvers, ok := objResolver.(map[string]interface{})[f.GoFieldName]; ok {
				fieldResolverType := reflect.TypeOf(objResolver)
				callbackType := reflect.TypeOf(fieldResolvers)
				arguments := []string{"ctx", "obj"}
				for _, arg := range f.Args {
					arguments = append(arguments, arg.VarName)
				}
				
				implementation = fmt.Sprintf(
					strings.Join(
						[]string{
							"    bundle, _ := r.container.bundles.%s()",
							"    resolvers := bundle.GraphqlResolver()",
							"    objectResolver := resolvers[\"%s\"].(%s)",
							"    callback := objectResolver[\"%s\"].(%s)",
							"",
							"    return callback(%s)",
						},
						"\n",
					),
					this.container.BundlePath(bundle),
					o.Name,
					fieldResolverType,
					f.GoFieldName,
					callbackType,
					strings.Join(arguments, ", "),
				)
			}
		}
	}
	
	return implementation
}
