package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	//"github.com/neo4j/neo4j-go-driver/v4/neo4j/internal/router"
)

type muxRouter struct{}

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

func (*muxRouter) GetParam(r *http.Request, param string) string {
	vars := mux.Vars(r)
	thisParam, ok := vars[param]
	if !ok {
		//err := "The parameter " + param + " does not exist"
		err := ERROR_MISSING_OR_NOT_FOUND_PARAMETER
		return err
	}
	return thisParam
}
