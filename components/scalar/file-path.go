package scalar

import (
	"strings"

	"bean/components/module/migrate"
)

type FilePath string

func (fp FilePath) String() string {
	out := string(fp)
	if strings.HasPrefix(out, "/") {
		return out
	}

	return migrate.RootDirectory() + "/" + out
}
