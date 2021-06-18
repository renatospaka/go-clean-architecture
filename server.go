package main

import (
	"fmt"
	"log"
	"net/http"

	router "github.com/renatospaka/go-clean-architecture/application/http"
)

var (
	//personRepo repository.PersonRepository = repository.NewPersonInMemoryRepository()
	httpRouter router.Router = router.NewMuxRouter()
)

func main() {
	const port = ":8000"
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("EMR Server running on port %s", port)
		fmt.Fprintln(w, "EMR Server running on port ", port)
	})

	httpRouter.SERVE(port)
}
