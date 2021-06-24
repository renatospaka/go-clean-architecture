package http

import "net/http"

type Router interface {
	GET(uri string, f func(w http.ResponseWriter, r *http.Request))
	POST(uri string, f func(w http.ResponseWriter, r *http.Request))
	SERVE(port string)
	GetParam(r *http.Request, param string) string
}

const (
	ERROR_MISSING_OR_NOT_FOUND_PARAMETER string = "parameter provided does not exist"
)
