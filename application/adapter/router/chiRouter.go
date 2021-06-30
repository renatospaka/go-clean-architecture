package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	//router "github.com/renatospaka/go-clean-architecture/infrastructure/http"
)

type chiRouter struct {}

var chiDispatcher = chi.NewRouter()

func NewChiRouter() Router {
	return &chiRouter{}
}

func (*chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, f)
}

func (*chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, f)
}

func (*chiRouter) SERVE(port string) {
	log.Printf("Chi http listening on port %s", port)
	log.Println(http.ListenAndServe(port, chiDispatcher))
}

func (*chiRouter) GetParam(r *http.Request, param string) string {
	panic("Not implemented yet")
}
