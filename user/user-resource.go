package main

import (
	"github.com/emicklei/go-restful"

	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// POST http://localhost:8080/users
// { "Id": "1", "Name": "Simi Pro"  }
//
// GET http://localhost:8080/users/1
//
// PUT http://localhost:8080/users/1
//{ "Id": "1", "Name": "Simi Pro"  }
//
// DELETE http://localhost:8080/users/1
//

type User struct {
	Id       string
	Name     string
	Email    string
	Password string
}

type UserResource struct {
	// TODO: Use Dao
	users map[string]User
}

func (u User) String() string {
	return fmt.Sprintf("Id: %s, Name: %s, Email: %s, Password: %s", u.Id, u.Name, u.Email, u.Password)
}

func (u UserResource) basicAuthentication(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authHeader := req.Request.Header.Get("Authorization")
	if len(authHeader) == 0 {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized, Im Sorry VIP Only")
		return
	}
	encoded := strings.SplitN(authHeader, " ", 2)

	if len(encoded) != 2 || encoded[0] != "Basic" {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized, Im Sorry VIP Only")
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(encoded[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized, Im Sorry VIP Only")
		return
	}
	authOK, user := u.Validate(pair[0], pair[1])
	if !authOK {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized, Im Sorry VIP Only")
		return
	}

	req.SetAttribute("user", user)
	chain.ProcessFilter(req, resp)

}

func (u UserResource) Validate(email, password string) (bool, User) {
	user := u.GetUserByEmail(email)
	if user.Password == password {
		return true, user
	}
	return false, User{}
}

func (u UserResource) GetUserByEmail(email string) User {
	for _, user := range u.users {
		if user.Email == email {
			return user
		}
	}
	return User{}
}

func (u UserResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/users").
		Doc("Manage Users").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{}))

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Returns(409, "duplicate user-id", nil).
		Reads(User{})) // from the request

	ws.Route(ws.POST("").To(u.createUser).
		// docs
		Doc("create a user").
		Operation("createUser").
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Operation("removeUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	container.Add(ws)

}

// GET http://localhost:8080/users/1
//
func (u UserResource) findUser(request *restful.Request, response *restful.Response) {
	log.Println("Request from User: ", request.Attribute("user"))
	id := request.PathParameter("user-id")
	usr := u.users[id]
	if len(usr.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
		return
	}
	response.WriteEntity(usr)
}

// POST http://localhost:8080/users
// { "Id": "1", "Name": "Simi Pro"  }
//
func (u *UserResource) createUser(request *restful.Request, response *restful.Response) {
	usr := &User{}
	err := request.ReadEntity(usr)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	usr.Id = strconv.Itoa(len(u.users) + 1) // simple id generation
	u.users[usr.Id] = *usr
	response.WriteHeaderAndEntity(http.StatusCreated, usr)
}

// PUT http://localhost:8080/users/1
// { "Id": "1", "Name": "Simi Pro"  }
//
func (u *UserResource) updateUser(request *restful.Request, response *restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	u.users[usr.Id] = *usr
	response.WriteEntity(usr)
}

// DELETE http://localhost:8080/users/1
//
func (u *UserResource) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	delete(u.users, id)
}

func main() {
	restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))
	wsContainer := restful.NewContainer()
	u := UserResource{map[string]User{}}
	u.users["1"] = User{Id: "1", Name: "Simi", Email: "simi", Password: "pro"} // default user

	wsContainer.Filter(u.basicAuthentication)
	u.Register(wsContainer)

	// we bind to 3001 which is the proxy port of gin
	log.Printf("start listening on localhost:3000 ")
	server := &http.Server{Addr: ":3001", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
