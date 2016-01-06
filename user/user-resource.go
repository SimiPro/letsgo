package main

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"github.com/letsgoli/common"

//	"github.com/letsgoli/images"

	"log"
	"net/http"
	appContext "golang.org/x/net/context"
	"github.com/letsgoli/Daos"
	"github.com/SimiPro/alice"
	"google.golang.org/appengine"
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
type UserResource struct {
	// TODO: Use Dao
	users   map[string]common.User
	userDao *daos.UserDao
}

func (u UserResource) GetUserByEmailOrUsername(email string) (bool, common.User) {
	for _, user := range u.users {
		if user.Email == email || user.Username == email {
			return true, user
		}
	}
	return false, common.User{}
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
	log.Println("HELLO HELLO HELOO")
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
	print("HALLO HALLO HALLOO!!!! ")
	var usr common.User
	if err := c.BindJSON(&usr); err != nil {
		Abort500WithError(c, err.Error())
		return
	}
	//usr.Id = strconv.Itoa(len(u.users) + 1) // simple id generation
	//u.users[usr.Id] = usr
	u.userDao.AddUser(&usr)
	c.JSON(http.StatusCreated, usr)
}


func (u *UserResource) signup(w http.ResponseWriter, r *http.Request) {
}


func Abort500WithError(c *gin.Context, errorli string) {
	c.Header("Content-Type", "text/plain")
	c.AbortWithError(http.StatusInternalServerError, &Error500{err: errorli})
	log.Fatalf(errorli)
}

// PUT http://localhost:8080/users/1
// { "Id": "1", "Name": "Simi Pro"  }
//
func (u *UserResource) updateUser(c *gin.Context) {
	var user common.User
	if err := c.BindJSON(&user); err != nil {
		Abort500WithError(c, err.Error())
		return
	}
	log.Println("USER: %v", user)
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://holidayers-1180.appspot.com")
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

func (u *UserResource) Register(usergroup *gin.RouterGroup) {
	usergroup.OPTIONS("login")
	usergroup.GET("login", u.login)
	usergroup.GET("find/:user-id", u.findUser)
	usergroup.PUT("update/:user-id", u.updateUser)
	usergroup.DELETE("remove/:user-id", u.removeUser)
}

func init() {
	// init datastore
	ctx := appContext.Background()
	// "user db"
	u := UserResource{
		map[string]common.User{},
		daos.NewUserDao(ctx),
	}

	u.users["1"] = common.User{Id: "1", Username: "Simi", Image: "no profile image", Email: "simi", Password: "pro" } // default user
	// end user db

	h := alice.New(common.NewLoggingHandler, middleware, common.NewRecoverHandler).ThenFuncWithContext(ctx, handler)
	commonHandlers := alice.New(setAppEngineContext, common.NewLoggingHandler, common.NewRecoverHandler)
	http.Handle("/user/test", commonHandlers.Append(common.NewBasicAuthentication).ThenFuncWithContext(ctx, testHandler))
	//	http.Handle("/user/signup", commonHandlers.ThenFunc(signupHandler))

	http.Handle("/user/testli", h)
}

func setAppEngineContext(next alice.ContextHandler) alice.ContextHandler {
	fn := func(empty appContext.Context, w http.ResponseWriter, r *http.Request) {
		// we hang on our empty parent context the appengine context for further middlewares
		newContext := appengine.WithContext(empty, r)
		next.ServeHTTPContext(newContext, w, r)
	}
	return alice.ContextHandlerFunc(fn)
}


func handler(ctx appContext.Context, rw http.ResponseWriter, req *http.Request) {
	reqID := requestIDFromContext(ctx)
	rw.Write([]byte("Hello request ID %s\n" + reqID))
}




/*
// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, u *common.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*common.User, bool) {
	u, ok := ctx.Value(userKey).(*common.User)
	return u, ok
}
*/

type key int
const requestIDKey key = 0

func newContextWithRequestID(ctx appContext.Context, req *http.Request) appContext.Context {
	return appContext.WithValue(ctx, requestIDKey, req.Header.Get("X-Request-ID"))
}

func requestIDFromContext(ctx appContext.Context) string {
	return ctx.Value(requestIDKey).(string)
}


func middleware(h alice.ContextHandler) alice.ContextHandler {
	return alice.ContextHandlerFunc(func(ctx appContext.Context, rw http.ResponseWriter, req *http.Request) {
		ctx = newContextWithRequestID(ctx, req)
		h.ServeHTTPContext(ctx, rw, req)
	})
}


func signupHandler(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "POST") {
		w.Write([]byte("Holy Moly"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ClearHandler wraps an http.Handler and clears request values at the end
// of a request lifetime.

func testHandler(ctx appContext.Context, w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		w.Write([]byte("Holy Moly"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
