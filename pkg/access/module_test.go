package access

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"bean/pkg/access/api/fixtures"
	"bean/pkg/namespace"
	fNamespace "bean/pkg/namespace/api/fixtures"
	"bean/pkg/user"
	fUser "bean/pkg/user/api/fixtures"
	"bean/pkg/util"
)

func module() *AccessModule {
	db := util.MockDatabase()

	logger := util.MockLogger()
	id := util.MockIdentifier()
	mUser := user.NewUserModule(db, logger, id)
	mNamespace := namespace.NewNamespaceModule(db, logger, id)
	module := NewAccessModule(db, id, logger, mUser, mNamespace)
	util.MockInstall(mUser, db)
	util.MockInstall(mNamespace, db)
	util.MockInstall(module, db)

	return module
}

func Test_Create(t *testing.T) {
	ctx := context.Background()
	ass := assert.New(t)
	this := module()

	// create user
	iUser := fUser.NewUserCreateInputFixture()
	oUser, err := this.userModule.UserCreate(ctx, iUser)
	ass.NoError(err)

	// create namespace
	iNamespace := fNamespace.NamespaceCreateInputFixture()
	oNamespace, err := this.namespaceModule.NamespaceCreate(ctx, iNamespace)
	ass.NoError(err)

	// base input
	input := fixtures.SessionCreateInputFixture(oNamespace.Namespace.ID, "", "")

	// create namespace
	t.Run("test email not found", func(t *testing.T) {
		ass.True(true)

		fmt.Println("WIP", oUser.User.ID, oNamespace.Namespace.ID, input)
	})

	t.Run("test email inactive", func(t *testing.T) {
	})

	t.Run("test email unverified", func(t *testing.T) {

	})

	t.Run("test password unmatched", func(t *testing.T) {
	})

	t.Run("test membership not found", func(t *testing.T) {
	})
}

func Test_Query(t *testing.T) {
}
