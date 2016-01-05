package common
import (
	"strings"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"log"
	"google.golang.org/cloud/datastore"
	"golang.org/x/net/context"
)


// Instances a Basic Authentication Middleware which checks auth with db and sets user on gin.Context with key "user"
func BasicAuthentication(ctx context.Context, client datastore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			NotAuthorized(c)
			return
		}
		encoded := strings.SplitN(authHeader, " ", 2)

		if len(encoded) != 2 || encoded[0] != "Basic" {
			NotAuthorized(c)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(encoded[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || len(pair[0]) == 0 || len(pair[1]) == 0 {
			NotAuthorized(c)
			return
		}
		authOK, user := Validate(pair[0], pair[1], client, ctx)
		if !authOK {
			NotAuthorized(c)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func Validate(email, password string, client datastore.Client, ctx context.Context) (bool, User) {
	query := datastore.NewQuery("User").
		Filter("Email =", email).
		Filter("Password =", password).Order("-")
	var users []User

	if _, err := client.GetAll(ctx, query , &users); err != nil {
		log.Fatalf("NOOOOO QUERY Not worked :(")
	}
	if (len(users) == 1) {
		return true, users[0]
	}
	return false, User{}
}

func NotAuthorized(c *gin.Context) {
	// c.Header("WWW-Authenticate", "Basic realm=Protected Area") if you uncomment this line a window with password & username pops up on client site
	c.AbortWithStatus(401)
}