package repository

import (
	"errors"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	uuid "github.com/satori/go.uuid"

	"github.com/renatospaka/go-clean-architecture/core/entity/person"
	"github.com/renatospaka/go-clean-architecture/infrastructure/db"
	"github.com/renatospaka/go-clean-architecture/infrastructure/utils"
)

type personNeo4j struct{}

var (
	personNeo4jDb db.Neo4jSession = db.NewNeo4jSession()
)

func NewPersonNeo4jRepository() PersonRepository {
	return &personNeo4j{}
}

func (*personNeo4j) FindById(id string) (*entity.Person, error) {
	//personNeo4jDb = db.NewNeo4jSession()
	if err := personNeo4jDb.Connect(db.AccessRead); err != nil {
		log.Printf("person.FindById.Connect.error: %v", err)
		return &entity.Person{}, err
	}

	if err := personNeo4jDb.IsValid(); err != nil {
		log.Printf("person.FindById.IsValid.error: %v", err)
		return &entity.Person{}, err
	}
	defer personNeo4jDb.Session.Close()
	
	tx, err := personNeo4jDb.Session.BeginTransaction()
	if err != nil {
		log.Printf("person.FindById.BeginTransaction.error: %v", err)
		return &entity.Person{}, err
	}
	defer tx.Rollback()
	
	cypher := "MATCH (one:Person {uuid: $uuid}) " +
						"WITH one " +
						"MATCH (one)-[:BIRTH]->(dob:Day) " +
						"RETURN one.uuid As uuid, one.name AS name, one.middleName AS middleName, one.lastName AS lastName, one.gender AS gender, one.email AS email, dob.date AS dob "
	params := map[string]interface{}{
		"uuid": id,
	}
	result, err := tx.Run(cypher, params)
	if err != nil {
		log.Printf("person.FindById.Run.error: %v", err)
		return &entity.Person{}, err
	}
	
	found := false
	thisGuy := entity.Person{}
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
		//here is a little tricky
		dob, ok := record.Get("dob")
		if dob != nil && ok {
			dob2 = dob.(dbtype.Date)
		}
		if !ok || dob == nil { 
			dob2 = dbtype.Date{} 
		}

		found = true
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

	if thisGuy.ID == "" || !found {
		err := errors.New(entity.ERROR_PERSON_INVALID_ID)
		log.Printf("person.FindById.error: %v", err)
		return &entity.Person{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("person.FindById.Commit.error: %v", err)
		return &entity.Person{}, err
	}

	err = tx.Close()
	if err != nil {
		log.Printf("person.FindById.Close.error: %v", err)
		return &entity.Person{}, err
	}

	personNeo4jDb.Session.Close()
	return &thisGuy, nil
}


//Add a new person and link him/her to the proper family
func (*personNeo4j) Add(person *entity.Person) (*entity.Person, error) {
	err := person.IsValid()
	if err != nil {
		log.Printf("person.Add.error: %v", err)
		return &entity.Person{}, err
	}

	//personNeo4jDb = db.NewNeo4jSession()
	if err := personNeo4jDb.Connect(db.AccessWrite); err != nil {
		log.Printf("person.Add.Connect.error: %v", err)
		return &entity.Person{}, err
	}

	if err := personNeo4jDb.IsValid(); err != nil {
		log.Printf("person.Add.IsValid.error: %v", err)
		return &entity.Person{}, err
	}
	defer personNeo4jDb.Session.Close()
	
	tx, err := personNeo4jDb.Session.BeginTransaction()
	if err != nil {
		log.Printf("person.Add.BeginTransaction.error: %v", err)
		return &entity.Person{}, err
	}
	defer tx.Rollback()

	action := "updated"
	if person.ID == "" {
		action = "created"
		person.ID = uuid.NewV4().String()
	}
	
	//remove previous day of birth from this user 
	//until I figure out a better way to do this
  err = resolveBirthDay(person.ID, tx)
	if err != nil {
		log.Printf("person.Add.resolveBirthDay.error: %v", err)
		log.Println("person.Add.resolveBirthDay.error: entire transaction canceled")
		return &entity.Person{}, err
	}
	
	//execute just the adition or update of the person
  err = addPerson(person, tx)
	if err != nil {
		log.Printf("person.Add.addPerson.error: %v", err)
		log.Println("person.Add.addPerson.error: entire transaction canceled")
		return &entity.Person{}, err
	}

	//if everything is OK, then commit at once
	err = tx.Commit()
	if err != nil {
		log.Printf("person.Add.Commit.error: %v", err)
		log.Println("person.Add.Commit.error: entire transaction canceled")
		return &entity.Person{}, err
	}

	err = tx.Close()
	if err != nil {
		log.Printf("person.Add.Close.error: %v", err)
		return &entity.Person{}, err
	}
	
	log.Printf("person.Add: id %v %v for person %v %v", person.ID, action, person.Name, person.LastName)
	personNeo4jDb.Session.Close()
	return person, nil
}

func addPerson(person *entity.Person, tx neo4j.Transaction) error {
	person.Responsible = entity.Someone
	person.CreatedAt = time.Now()
	person.UpdatedAt = time.Now()

	cypher := "MERGE (one:Person {uuid: $uuid}) " +
						"ON CREATE SET one.name = $name, " +
								"one.fullName = TOLOWER(REPLACE($name, ' ', '') + REPLACE($middleName, ' ', '') + REPLACE($lastName, ' ', '')), " +
								"one.gender = $gender, " +
								"one.email = $email, " +
								"one.middleName = $middleName, "  +
								"one.lastName = $lastName, " +
								"one.createdAt = $created, " +
								"one.updatedAt = $updated " +
						"ON MATCH SET one.gender = $gender, " +
								"one.name = $name, " +
								"one.email = $email, " +
								"one.middleName = $middleName, " +
								"one.lastName = $lastName, " +
								"one.updatedAt = $updated " +  
						//IDENTIFY THE DAY OF BIRTH
						"WITH one " +
						"MATCH (dob:Day {uuid: $dob}) " +
						"MERGE (one)-[:BIRTH]->(dob) " +
						"SET one.ageInMonths = duration.between(date(dob.date), date()).months, " +
								"one.ageInYears = duration.between(date(dob.date), date()).years " +
						"RETURN one.uuid As uuid, one.name AS name, one.middleName AS middleName, one.lastName AS lastName, " +
								"one.gender AS gender, one.email AS email, dob.date AS dob, " +
								"one.ageInMonths AS ageInMonths, one.ageInYears AS ageInYears "
	params := map[string]interface{}{
		"uuid": person.ID,
		"name": person.Name,
		"middleName": person.MiddleName,
		"lastName": person.LastName,
		"gender": person.Gender,
		"email": person.Email,
		"dob": person.DOB.Format(utils.LayoutBRShort),
		"created": person.CreatedAt.Format(utils.LayoutBRShort),
		"updated": person.UpdatedAt.Format(utils.LayoutBRShort),
	}
	_, err := tx.Run(cypher, params)
	return err
}

func resolveBirthDay(id string, tx neo4j.Transaction) error {
	cypher := "MATCH (one:Person {uuid: $uuid}) " +
						"MATCH (one)-[bday:BIRTH]->(:Day) " +
						"DELETE bday" 
	params := map[string]interface{}{
		"uuid": id,
	}
	_, err := tx.Run(cypher, params)
	return err
}
