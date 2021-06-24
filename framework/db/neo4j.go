package db

import (
	"errors"
	//"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const (
	ERROR_DB_MISSING_URI         string = "cannot reach database's URI"
	ERROR_DB_CREDENTIALS         string = "user or password incorrect"
	ERROR_DB_CONNECTION          string = "unable to connect to database server"
	ERROR_DB_MISSING_CONFIG_FILE string = "unable to locate .env file"
	ERROR_DB_NO_CONNECTION       string = "connection to database missing"
)

type Neo4jSession struct {
	Session neo4j.Session
	driver  neo4j.Driver
}


func NewNeo4jSession(accessMode neo4j.AccessMode) Neo4jSession {
	var session neo4j.Session

	driverSession, err := connect()
	if err != nil {
		return Neo4jSession{}
	}
	session = driverSession.NewSession(neo4j.SessionConfig{AccessMode: accessMode})

	return Neo4jSession{
		driver:  driverSession,
		Session: session,
	}
}


func connect() (neo4j.Driver, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New(ERROR_DB_MISSING_CONFIG_FILE)
	}

	uri, found := os.LookupEnv("NEO4J_URI")
	if !found {
		return nil, errors.New(ERROR_DB_MISSING_URI)
	}

	username, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		return nil, errors.New(ERROR_DB_CREDENTIALS)
	}

	password, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		return nil, errors.New(ERROR_DB_CREDENTIALS)
	}

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		// fmt.Println(err.Error())
		return nil, errors.New(ERROR_DB_CONNECTION)
	}

	return driver, nil
}


func (d *Neo4jSession) isValid() error {
	if d == nil {
		return errors.New(ERROR_DB_NO_CONNECTION)
	}

	return nil
}
