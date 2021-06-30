package db

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jSession struct {
	Session neo4j.Session
	driver  neo4j.Driver
}

// type neo4jDb struct{
// 	db Neo4jSession
// }

var (
	session neo4j.Session
	driver  neo4j.Driver
)

func connect(accessType int) (neo4j.Driver, neo4j.Session, error) {
	var (
		sessionConfig  neo4j.SessionConfig
		accessTypeDesc string
	)

	problem := false
	problemMessage := ""
	err := godotenv.Load()
	if err != nil {
		log.Printf("neo4j: %v", ERROR_DB_MISSING_CONFIG_FILE)
		problem = true
		problemMessage = ERROR_DB_MISSING_CONFIG_FILE
		//return nil, nil, errors.New(ERROR_DB_MISSING_CONFIG_FILE)
	}

	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	if uri == "" {
		log.Printf("neo4j: %v", ERROR_DB_MISSING_URI)
		problem = true
		problemMessage = ERROR_DB_MISSING_URI
		//return nil, nil, errors.New(ERROR_DB_MISSING_URI)
	}

	if username == "" {
		log.Printf("neo4j: %v", ERROR_DB_CREDENTIALS)
		problem = true
		problemMessage = ERROR_DB_CREDENTIALS
		//return nil, nil, errors.New(ERROR_DB_CREDENTIALS)
	}

	if password == "" {
		log.Printf("neo4j: %v", ERROR_DB_CREDENTIALS)
		problem = true
		problemMessage = ERROR_DB_CREDENTIALS
		//return nil, nil, errors.New(ERROR_DB_CREDENTIALS)
	}

	if problem {
		return nil, nil, errors.New(problemMessage)
	}

	driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Printf("neo4j: %v", ERROR_DB_UNABLE_2_CONNECT)
		return nil, nil, errors.New(ERROR_DB_UNABLE_2_CONNECT)
	}
	
	err = driver.VerifyConnectivity()
	if err != nil {
		log.Printf("neo4j: %v", err)
		return nil, nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	log.Println("neo4j: driver connected to the server")

	if accessType == AccessWrite {
		accessTypeDesc = "write"
		sessionConfig.AccessMode = neo4j.AccessModeWrite
	} else {
		accessTypeDesc = "read"
		sessionConfig.AccessMode = neo4j.AccessModeRead
	}
	session = driver.NewSession(sessionConfig)
	
	log.Printf("neo4j: session established in %v mode", accessTypeDesc)
	
	return driver, session, nil
}

func NewNeo4jSession() Neo4jSession {
	return Neo4jSession{}
}

//Open a session and connect to the server
//According to CQRS, you may choose between
//	accessType =
//		AccessRead
//		AccessWrite
//to better performance
func (d *Neo4jSession) Connect(accessType int) error {
	driver, session, err := connect(accessType)
	if err != nil {
		log.Printf("neo4j: %v", err.Error())
		
		return err
	}

	// err = d.IsValid()
	// if err != nil {
	// 	return err
	// }
	
	d.Session = session
	d.driver = driver
	return nil
}

//Check whenever the session is valid and connected to the server
func (d *Neo4jSession) IsValid() error {
	//if driver wasn't even created before
	if d == nil {
		return errors.New(ERROR_DB_MISSING_CONNECTION)
	}
	
	//if driver is set but there are errors
	err := d.driver.VerifyConnectivity()
	if err != nil {
		log.Printf("neo4j.IsValid: %v", ERROR_DB_MISSING_CONNECTION)
		return errors.New(ERROR_DB_MISSING_CONNECTION)
	}
	log.Println("neo4j: session is valid")

	return nil
}

// func handleError(err error) bool {
// 	var e *neo4j.Neo4jError
	
// 	isErr := neo4j.IsNeo4jError(err)
// 	if isErr {
// 		e = err.(*neo4j.Neo4jError)
// 		if e.
// 	}

// 	return true
// }

// func (d *Neo4jSession) HanddleError(err error) {
// 	panic("not implemented")
// }