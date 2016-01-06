package common
import (
	"net/http"
	"log"
	"github.com/SimiPro/alice"
	"golang.org/x/net/context"
)

func NewRecoverHandler(next alice.ContextHandler) alice.ContextHandler {
	fn := func(c context.Context, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTPContext(c, w, r)
	}
	return alice.ContextHandlerFunc(fn)
}