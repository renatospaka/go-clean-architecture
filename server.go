package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"github.com/renatospaka/go-clean-architecture/application/controller"
	router "github.com/renatospaka/go-clean-architecture/application/http"
	"github.com/renatospaka/go-clean-architecture/framework/db"
)

var (
	//personRepo repository.PersonRepository = repository.NewPersonInMemoryRepository()
	httpRouter router.Router = router.NewMuxRouter()
	personController controller.PersonController = controller.NewPersonController()
	dbServer db.Neo4jSession = db.NewNeo4jSession(neo4j.AccessModeWrite)
)

func main() {
	//database server
	err := dbServer.IsValid()
	if err != nil {
		log.Println("Could not connect to database server due to", err.Error())
	}

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
