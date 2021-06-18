package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct {}

var muxDispatcher = mux.NewRouter()

func NewMuxRouter() Router {
	return &muxRouter{}
}

func (*muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("GET")
}

func (*muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, f).Methods("POST")
}

func (*muxRouter) SERVE(port string) {
	log.Printf("Mux http listening on port %s", port)
	log.Println(http.ListenAndServe(port, muxDispatcher))
}
