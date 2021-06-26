package repository

import (
	"errors"
	"log"
	"time"

	//_ "github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/framework/db"
	"github.com/renatospaka/go-clean-architecture/framework/utils"
	uuid "github.com/satori/go.uuid"
)

type personNeo4j struct{}

var (
	personNeo4jDbArray []entity.Person
	//personNeo4jDb db.Neo4jSession = db.NewNeo4jSessionWrite()
)

func NewPersonNeo4jRepository() PersonRepository {
	return &personNeo4j{}
}

func (*personNeo4j) FindById(id string) (*entity.Person, error) {
	personNeo4jDb := db.NewNeo4jSessionWrite()
	err := personNeo4jDb.IsValid()
	if err != nil {
		return &entity.Person{}, err
	}
	defer personNeo4jDb.Session.Close()

	tx, err := personNeo4jDb.Session.BeginTransaction()
	if err != nil {
		log.Printf("person.findbyid.BeginTransaction.error: %v", err)
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("person.findbyid.result.Single.Rollback.error: %v", err2)
		}
		return &entity.Person{}, err
	}
	defer tx.Close()
	
	//cypher := "CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)"
	cypher := "MATCH (p:Person {uuid: $uuid}) " +
						"RETURN p.uuid As uuid, p.name AS name, p.middleName AS middleName, p.lastName AS lastName, p.gender AS gender, p.email AS email"
	params := map[string]interface{}{
									"uuid": id,
								}
	result, err := tx.Run(cypher, params)
	if err != nil {
		log.Printf("person.findbyid.Run.error: %v", err)
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("person.findbyid.result.Single.Rollback.error: %v", err2)
		}
		return &entity.Person{}, err
	}
	
	// if result.Next() {
	record, err := result.Single()
	if err != nil {
		log.Printf("person.findbyid.result.Single.error: %v", err)
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("person.findbyid.result.Single.Rollback.error: %v", err2)
		}
		return &entity.Person{}, err
	}

	thisGuy := entity.Person{
		ID: utils.IsNilString(record.Values[0]),
		Name: utils.IsNilString(record.Values[1]),
		MiddleName: utils.IsNilString(record.Values[2]),
		LastName: utils.IsNilString(record.Values[3]),
		Email: utils.IsNilString(record.Values[4]),
		Gender: utils.IsNilString(record.Values[5]),
		DOB: utils.IsNilTime(time.Time{}),
	}

	if thisGuy.ID == "" {
		err := errors.New(entity.ERROR_PERSON_INVALID_ID)
		log.Printf("person.findbyid.error: %v", err)
		err2 := tx.Rollback()
		if err2 != nil {
			log.Printf("person.findbyid.result.Single.Rollback.error: %v", err2)
		}
		return &entity.Person{}, err
	}

	tx.Commit()
	tx.Close()
	personNeo4jDb.Session.Close()
	return &thisGuy, nil
}

func (*personNeo4j) Add(newPerson *entity.Person) (*entity.Person, error) {
	err := newPerson.IsValid()
	if err != nil {
		log.Printf("person.add.error: %v", err)
		return &entity.Person{}, err
	}

	newPerson.ID = uuid.NewV4().String()
	newPerson.Responsible = entity.Someone
	newPerson.CreatedAt = time.Now()
	newPerson.UpdatedAt = time.Now()

	return newPerson, nil
}
