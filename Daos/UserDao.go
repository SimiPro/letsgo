package daos
import (
	"google.golang.org/cloud/datastore"
	"golang.org/x/net/context"
	"github.com/letsgoli/common"
	"log"
)

type UserDao struct {
	client datastore.Client
	ctx context.Context
}

func NewDao(_client datastore.Client, _ctx context.Context) {
	return &UserDao{
		client: _client,
		ctx: _ctx,
	}
}

/**
	Returns first User with this email
 */
func (u *UserDao) GetUserByEmail(email string){
	query := datastore.NewQuery("User").
	Filter("Email =", email)

	var users []common.User
	if _, err := u.client.GetAll(u.ctx, query , &users); err != nil {
		log.Fatalf("NOOOOO QUERY Not worked :(")
	}
	if (len(users) == 1) {
		return true, users[0]
	}
	return nil
}

func (u *UserDao) AddUser(user common.User) {

}
