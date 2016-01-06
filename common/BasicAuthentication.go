package common

import (
	"strings"
	"encoding/base64"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/http"
	"github.com/SimiPro/alice"
)

type BasicAuthentication struct {
	next alice.ContextHandler
}

func NewBasicAuthentication(next alice.ContextHandler) alice.ContextHandler {
	return BasicAuthentication{next: next}
}

func (b BasicAuthentication) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Instances a Basic Authentication Middleware which checks auth with db and sets user on gin.Context with key "user"
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		NotAuthorized(w)
		return
	}
	encoded := strings.SplitN(authHeader, " ", 2)
	if len(encoded) != 2 || encoded[0] != "Basic" {
		NotAuthorized(w)
		return
	}
	payload, _ := base64.StdEncoding.DecodeString(encoded[1])
	pair := strings.SplitN(string(payload), ":", 2)
	log.Printf("1: %v, 2: %v", pair[0], pair[1])
	if len(pair) != 2 || len(pair[0]) == 0 || len(pair[1]) == 0 {
		NotAuthorized(w)
		return
	}
	authOK, user := Validate(pair[0], pair[1], ctx)
	if !authOK {
		NotAuthorized(w)
		return
	}
	//.Value(r, "user", user)
	newContext := UserNewContext(ctx, &user)
	b.next.ServeHTTPContext(newContext, w, r)
}

func Validate(email, password string, ctx context.Context) (bool, User) {
	query := datastore.NewQuery("User").
	Filter("Email =", email).
	Filter("Password =", password)

	result := query.Run(ctx)
	var user User
	if _, err := result.Next(&user); err != nil {
		log.Println("No User found")
		return false, User{}
	}
	log.Println("USER FOUND!")
	return true, user
}

func NotAuthorized(w http.ResponseWriter) {
	// c.Header("WWW-Authenticate", "Basic realm=Protected Area") if you uncomment this line a window with password & username pops up on client site
	w.WriteHeader(401)
}