package util

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
)

func ParseFile(path string) ([]byte, error) {
	raw, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, err
	}

	return replaceEnvVariables(raw)
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
