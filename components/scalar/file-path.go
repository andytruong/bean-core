package scalar

import (
	"strings"

	"bean/pkg/util/migrate"
)

type FilePath string

func (this FilePath) String() string {
	out := string(this)
	if strings.HasPrefix(out, "/") {
		return out
	}

	return migrate.RootDirectory() + "/" + out
}
