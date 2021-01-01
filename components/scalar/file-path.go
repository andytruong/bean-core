package scalar

import (
	"path"
	"runtime"
	"strings"
)

type FilePath string

func (fp FilePath) String() string {
	out := string(fp)
	if strings.HasPrefix(out, "/") {
		return out
	}

	return RootDirectory() + "/" + out
}

func RootDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Dir(filename)
	dir = path.Dir(dir)
	dir = path.Dir(dir)

	return dir
}
