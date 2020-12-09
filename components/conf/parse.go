package conf

import (
	"bytes"
	"errors"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

var ErrExceededRefLimit = errors.New("exceeded limit of nested references")
var refLimit = 1000

func ParseFile(path string, out interface{}) error {
	content, err := read(path)

	if nil != err {
		return err
	}

	content, err = replaceEnvVariables(content)
	if nil != err {
		return err
	}

	err = yaml.Unmarshal(content, out)

	return err
}

func replaceEnvVariables(inBytes []byte) ([]byte, error) {
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
