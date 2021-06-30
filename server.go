package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/renatospaka/go-clean-architecture/application/controller"
	"github.com/renatospaka/go-clean-architecture/application/adapter/router"
)

var (
	httpRouter router.Router = router.NewMuxRouter()
	personController controller.PersonController = controller.NewPersonController()
)

func main() {
	//http server
	const port = ":8000"
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("EMR Server running on port %s", port)
		fmt.Fprintln(w, "EMR Server running on port ", port)
	})

	httpRouter.GET("/person/{id}", personController.GetPerson)
	httpRouter.POST("/person", personController.AddPerson)

	httpRouter.SERVE(port)
}
