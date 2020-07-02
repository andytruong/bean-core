package scalar

import (
	"strings"

	"bean/components/module/migrate"
)

type FilePath string

func (this FilePath) String() string {
	out := string(this)
	if strings.HasPrefix(out, "/") {
		return out
	}

	return migrate.RootDirectory() + "/" + out
}
