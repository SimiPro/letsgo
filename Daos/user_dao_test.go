package daos_test

import (
	"testing"

	"google.golang.org/appengine/aetest"
	"github.com/letsgoli/Daos"
	"github.com/letsgoli/common"
	"google.golang.org/appengine/datastore"
)

type Testli struct {
	name string
}

func Test_test(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	testlit := &Testli{
		name: "simipro",
	}

	key := datastore.NewKey(ctx, "Testlis", "", 1,nil)
	if _, err := datastore.Put(ctx, key, testlit); err != nil {
		t.Fatal(err)
	}


}

func Test_AddUser(t *testing.T) {
		ctx, done, err := aetest.NewContext()
		if err != nil {
			t.Fatal(err)
		}
		userDao := daos.NewUserDao(ctx)

		user := &common.User {
			Email: "Email",
			Username: "Username",
			Firstname: "Firstname",
			Lastname: "Lastname",
			Password:"Password",
		}
		userDao.AddUser(user)

		defer done()

		// Run code and tests requiring the context.Context using ctx.
}


