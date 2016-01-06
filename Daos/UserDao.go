package daos
import (
	"golang.org/x/net/context"
	"github.com/letsgoli/common"
	"log"
"google.golang.org/appengine/datastore"
)

type UserDao struct {
	ctx context.Context
}

func NewUserDao( _ctx context.Context) *UserDao {
	return &UserDao{
		ctx: _ctx,
	}
}

/**
	Returns first User with this email
 */
func (u *UserDao) GetUserByEmail(email string) *common.User {
	query := datastore.NewQuery("User").
	Filter("Email =", email)

	var users []common.User

	if _, err := query.GetAll(u.ctx, &users); err != nil {
		log.Fatalf("NOOOOO QUERY Not worked :(")
	}
	if (len(users) == 1) {
		return &users[0]
	}
	return nil
}

func (u *UserDao) AddUser(user *common.User) {
	key := datastore.NewIncompleteKey(u.ctx, "User", nil)
	log.Println("Key: %v", key)
	if _, err := datastore.Put(u.ctx, key, user); err != nil {
		log.Fatalf("Failed to insert User in db: %v" , err.Error())
	} else {
		log.Println("User inserted!")
	}
}
