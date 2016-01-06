package common

import (
	"net/http"
	"time"
	"log"
	"github.com/SimiPro/alice"
	"golang.org/x/net/context"
)

type LoggingHandler struct  {
	next alice.ContextHandler
}

// called with each request newly
func NewLoggingHandler(next alice.ContextHandler) alice.ContextHandler {
	return LoggingHandler{next: next}
}

func (l LoggingHandler) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	l.next.ServeHTTPContext(ctx, w, r)
	t2 := time.Now()
	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}