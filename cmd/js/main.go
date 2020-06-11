package main

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func main() {
	{
		loader := gojsonschema.NewSchemaLoader()
		source := gojsonschema.NewStringLoader(`{ "kind": "number" }`)

		_, err := loader.Compile(source)

		fmt.Println("error : ", err)
	}

	{
		loader := gojsonschema.NewSchemaLoader()
		loader.Validate = true
		err := loader.AddSchemas(
			gojsonschema.NewStringLoader(
				`{ "whattt": "string" }`,
			),
		)

		fmt.Println("error 2: ", err, loader.AutoDetect)
	}
}
