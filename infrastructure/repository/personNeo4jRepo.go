package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/renatospaka/go-clean-architecture/entity"
	"github.com/renatospaka/go-clean-architecture/framework/db"

	//"github.com/renatospaka/go-clean-architecture/framework/utils"
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
	thisGuy := entity.Person{}
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
	cypher := "MATCH (one:Person {uuid: $uuid}) " +
		"WITH one " +
		"MATCH (one)-[:BIRTH]->(dob:Day) " +
		"RETURN one.uuid As uuid, one.name AS name, one.middleName AS middleName, one.lastName AS lastName, one.gender AS gender, one.email AS email, dob.date AS dob "
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

	for result.Next() {
		var dob2 dbtype.Date
		record := result.Record()
		uuid, ok := record.Get("uuid")
		name, ok := record.Get("name")
		middleName, ok := record.Get("middleName")
		lastName, ok := record.Get("lastName")
		email, ok := record.Get("email")
		if !ok || email == nil {
			email = ""
		}
		gender, ok := record.Get("gender")
		if !ok || gender == nil {
			gender = ""
		}

		dob, ok := record.Get("dob")
		if dob != nil {
			dob2 = dob.(dbtype.Date)
		}
		if !ok || dob == nil {
			dob2 = dbtype.Date{}
		}

		thisGuy = entity.Person{
			ID:         uuid.(string),
			Name:       name.(string),
			MiddleName: middleName.(string),
			LastName:   lastName.(string),
			Email:      email.(string),
			Gender:     gender.(string),
			DOB:        time.Time(dob2),	// that gave me 2 entire days of a huge headache
		}
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
