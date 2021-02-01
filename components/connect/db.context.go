package connect

import (
	"context"

	"gorm.io/gorm"

	"bean/components/scalar"
)

const DatabaseContextKey scalar.ContextKey = "bean.db"

func WithContextValue(ctx context.Context, wrapper *Wrapper) context.Context {
	return context.WithValue(ctx, DatabaseContextKey, wrapper)
}

func DB(ctx context.Context) *gorm.DB {
	if db, ok := ctx.Value(DatabaseContextKey).(*gorm.DB); ok {
		return db
	}

	if wrapper, ok := ctx.Value(DatabaseContextKey).(*Wrapper); ok {
		db, err := wrapper.Master()
		if nil != err {
			panic(err)
		}

		return db
	}

	return nil
}

func DBToContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, DatabaseContextKey, db)
}
