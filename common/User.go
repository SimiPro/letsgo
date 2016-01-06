package common

import (
"fmt"
"golang.org/x/net/context"
)

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"  binding:"required"`
	Firstname string `json:"firstName" binding:"required"`
	Lastname  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Image     string `json:"image"`
	Password  string `json:"password"`
}

func (u User) String() string {
	return fmt.Sprintf("Id: %s, Name: %s, Email: %s, Password: %s", u.Id, u.Username, u.Email, u.Image, u.Password)
}

type key int
const userKey key = 17

// NewContext returns a new Context that carries value u.
func UserNewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

