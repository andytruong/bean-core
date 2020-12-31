package util

import (
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"bean/components/scalar"
)

func MockLogger() *zap.Logger {
	return zap.NewNop()
}

func MockIdentifier() *scalar.Identifier {
	return &scalar.Identifier{}
}
