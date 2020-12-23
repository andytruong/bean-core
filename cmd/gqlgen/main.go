package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin/resolvergen"
	"github.com/pkg/errors"

	"bean/components/util"
	"bean/pkg/infra"
)

func main() {
	log.SetOutput(ioutil.Discard)
	_ = os.Setenv("ENV", "test")
	path := "config.yaml"
	if path == "" {
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
		err = api.Generate(cfg, api.AddPlugin(plugin{container: container}))

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(3)
		}
	}
}

type plugin struct {
	container *infra.Container
}

func (p plugin) Name() string {
	return "bean"
}

func (p plugin) GenerateCode(data *codegen.Data) error {
	build := &resolvergen.ResolverBuild{
		File:         &resolvergen.File{},
		PackageName:  data.Config.Resolver.Package,
		ResolverType: data.Config.Resolver.Type,
		HasRoot:      true,
	}

	for _, o := range data.Objects {
		if o.HasResolvers() {
			build.File.Objects = append(build.File.Objects, o)
		}
		for _, f := range o.Fields {
			if !f.IsResolver {
				continue
			}

			build.File.Resolvers = append(build.File.Resolvers, p.newResolver(data, o, f))
		}
	}

	options := templates.Options{
		Filename:    "pkg/infra/graphql.resolvers.go",
		PackageName: "infra",
		FileNotice:  `// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.`,
		Data:        build,
		Packages:    data.Config.Packages,
		Funcs: map[string]interface{}{
			"resolverBody": func(input string, correctType string) string {
				return strings.Replace(input, "%tempType%", "func "+correctType, 1)
			},
		},
	}

	return templates.Render(options)
}

func (p plugin) newResolver(data *codegen.Data, o *codegen.Object, f *codegen.Field) *resolvergen.Resolver {
	resolver := &resolvergen.Resolver{
		Object:         o,
		Field:          f,
		Implementation: p.resolverBody(data, o, f),
	}

	return resolver
}

func (p plugin) resolverBody(data *codegen.Data, o *codegen.Object, f *codegen.Field) string {
	implementation := fmt.Sprintf(`panic("no implementation found in resolvers[%s][%s]")`, o.Name, f.GoFieldName)

	for _, bundle := range p.container.BundleList() {
		resolvers := bundle.GraphqlResolver()
		if nil == resolvers {
			continue
		}

		if objResolver, ok := resolvers[o.Name]; ok {
			if _, ok := objResolver.(map[string]interface{})[f.GoFieldName]; ok {
				fieldResolverType := reflect.TypeOf(objResolver)
				arguments := []string{"ctx"}

				if !f.Object.Root {
					arguments = append(arguments, "obj")
				}

				for _, arg := range f.Args {
					arguments = append(arguments, arg.VarName)
				}

				implementation = fmt.Sprintf(
					strings.Join(
						[]string{
							// TODO: Handle @requireAuth()
							"    bundle, _ := r.container.bundles.%s()",
							"    resolvers := bundle.GraphqlResolver()",
							"    objectResolver := resolvers[\"%s\"].(%s)",
							"    callback := objectResolver[\"%s\"].(%s)",
							"",
							"    return callback(%s)",
						},
						"\n",
					),
					p.container.BundlePath(bundle),
					o.Name,
					fieldResolverType,
					f.GoFieldName,
					"%tempType%",
					strings.Join(arguments, ", "),
				)
			}
		}
	}

	return implementation
}
