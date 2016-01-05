package common

import "fmt"

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