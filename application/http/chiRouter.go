package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type chiRouter struct {}

var chiDispatcher = chi.NewRouter()

func NewChiRouter() Router {
	return &chiRouter{}
}

func (*chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.HandleFunc(uri, f).Methods("GET")
}

func (*chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (*chiRouter) SERVE(port string) {
	log.Printf("Chi http listening on port %s", port)
	log.Println(http.ListenAndServe(port, chiDispatcher))
}
