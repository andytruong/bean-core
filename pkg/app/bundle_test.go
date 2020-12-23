package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"bean/components/util"
)

func bundle(ctx context.Context) *AppBundle {
	con := util.MockDatabase().WithContext(ctx)
	idr := util.MockIdentifier()
	log := util.MockLogger()
	bun, _ := NewApplicationBundle(con, idr, log)

	util.MockInstall(bun, con)

	return bun
}

func Test(t *testing.T) {
	ass := assert.New(t)
	ctx := context.Background()
	bun := bundle(ctx)

	ass.True(
		bun.con.Migrator().HasTable("applications"),
		"Table `applications` was created",
	)
}
