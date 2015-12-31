package main

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

//	"github.com/letsgoli/images"

	"encoding/base64"
	"fmt"
	"log"
	"net/http"
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
	Id       string `json:"id"`
	Username string `json:"name"  binding:"required"`
	Email    string `json:"email" binding:"required"`
	Image 	 string `json:"image" binding:"required"`
	Password string `json:"password"`
}

type UserResource struct {
	// TODO: Use Dao
	users map[string]User
}

func (u User) String() string {
	return fmt.Sprintf("Id: %s, Name: %s, Email: %s, Password: %s", u.Id, u.Username, u.Email, u.Image, u.Password)
}

func NotAuthorized(c *gin.Context) {
	c.Header("WWW-Authenticate", "Basic realm=Protected Area")
	c.AbortWithStatus(401)
}
func isOptionsCorsRequest(c *gin.Context) bool {
	ORIGIN_HEADER := "Origin"
	OPTIONS_METHOD := "OPTIONS"

	if c.Request.Method == OPTIONS_METHOD  {
		if c.Request.Header.Get(ORIGIN_HEADER) != "" {
			return true
		}
	}
	return false
}

func (u UserResource) basicAuthentication(c *gin.Context) {
	if isOptionsCorsRequest(c) {
		// check if "preflight" cors request if yes dont check authorization
		// TODO: REMOVE IN PRODUCTION
		c.Next()
	} else {
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
		authOK, user := u.Validate(pair[0], pair[1])
		if !authOK {
			NotAuthorized(c)
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func (u UserResource) Validate(email, password string) (bool, User) {
	found, user := u.GetUserByEmail(email)
	if found && user.Password == password {
		return true, user
	}
	return false, User{}
}

func (u UserResource) GetUserByEmail(email string) (bool, User) {
	for _, user := range u.users {
		if user.Email == email {
			return true, user
		}
	}
	return false, User{}
}

func (u UserResource) Register(group *gin.RouterGroup) {
	group.OPTIONS("login")
	group.GET("login", u.login)
	group.GET("find/:user-id", u.findUser)
	group.PUT("update/:user-id", u.updateUser)
	group.POST("create", u.createUser)
	group.DELETE("remove/:user-id", u.removeUser)
}

type ErrorUserNotFound struct {
}

func (err ErrorUserNotFound) Error() string {
	return "User Not Found"
}

type Error500 struct {
	err string
}

func (errorli Error500) Error() string {
	return errorli.err
}


// Get http://localhost:8080/login --> Returns MyUser
func (u UserResource) login(c *gin.Context) {
	c.JSON(http.StatusOK, c.MustGet("user"));
}

// GET http://localhost:8080/users/1
//
func (u UserResource) findUser(c *gin.Context) {
	log.Println("Request from User: ", c.MustGet("user"))
	id := c.Param("user-id")
	usr := u.users[id]
	if len(usr.Id) == 0 {
		c.Header("Content-Type", "text/plain")
		c.AbortWithError(http.StatusNotFound, new(ErrorUserNotFound))
		return
	}
	c.JSON(http.StatusOK, usr)
}

// POST http://localhost:8080/users
// { "Id": "1", "Name": "Simi Pro"  }
//
func (u *UserResource) createUser(c *gin.Context) {
	var usr User
	if err := c.BindJSON(&usr); err != nil {
		Abort500WithError(c, err.Error())
		return
	}
	usr.Id = strconv.Itoa(len(u.users) + 1) // simple id generation
	u.users[usr.Id] = usr
	c.JSON(http.StatusCreated, usr)
}

func Abort500WithError(c *gin.Context, errorli string) {
	c.Header("Content-Type", "text/plain")
	c.AbortWithError(http.StatusInternalServerError, &Error500{err: errorli})
}

// PUT http://localhost:8080/users/1
// { "Id": "1", "Name": "Simi Pro"  }
//
func (u *UserResource) updateUser(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		Abort500WithError(c, err.Error())
		return
	}
	u.users[user.Id] = user
	c.JSON(http.StatusOK, user)
}

// DELETE http://localhost:8080/users/1
//
func (u *UserResource) removeUser(c *gin.Context) {
	id := c.MustGet("user-id").(string)
	delete(u.users, id)
}

func AuthRequiered(c *gin.Context) {
	log.Println("Auth Middleware in action")
	c.Next()
}

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}

func CorsHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:63343")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
		c.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	// here could operations be executed which only must be executed once. e.g load a key file or so

	return func(c *gin.Context) {
		token := c.Request.FormValue("api_token")
		if token == "" {
			respondWithError(401, "API token required", c)
			return
		}
		if validateToken("myemail", token) {
			respondWithError(401, "Invalid token", c)
			return
		}
		c.Next()
	}
}

func validateToken(email, token string) bool {
	return true
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.AbortWithStatus(code)
}

func main() {

	// Gin Middleware config
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(RequestIdMiddleware())
	router.Use(CorsHeader())
	// end gin middleware config

	// "user db"
	u := UserResource{map[string]User{}}
	u.users["1"] = User{Id: "1", Username: "Simi", Image: "no profile image" , Email: "simi", Password: "pro"} // default user
	// end user db

	// url with "/user"
	userGroup := router.Group("/user")
	userGroup.Use(u.basicAuthentication)
	u.Register(userGroup)


	//url images
	imageGroup := router.Group("/image")
	imageGroup.Use(u.basicAuthentication)
	//	images.Register(imageGroup)

	// we bind to 3001 which is the proxy port of gin
	log.Printf("start listening on localhost:3000 ")
	router.Run(":8000")

	//server := &http.Server{Addr: ":3001", Handler: wsContainer}
	//log.Fatal(server.ListenAndServe())
}
