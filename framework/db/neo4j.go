package db

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Neo4jSession struct {
	Session neo4j.Session
	driver  neo4j.Driver
}

type db struct{}

func NewNeo4jSessionRead() Neo4jSession {
	var (
		session          neo4j.Session
		//thisNeo4jSession Neo4jSession = Neo4jSession{}
	)

	driverSession, err := connect()
	if err != nil {
		return Neo4jSession{}
	}
	
	session = driverSession.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	return Neo4jSession{
		Session: session,
		driver: driverSession,
	}
}

func NewNeo4jSessionWrite() Neo4jSession {
	var (
		session          neo4j.Session
		//thisNeo4jSession Neo4jSession = Neo4jSession{}
	)

	driverSession, err := connect()
	if err != nil {
		return Neo4jSession{}
	}
	
	session = driverSession.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	return Neo4jSession{
		Session: session,
		driver: driverSession,
	}
}

//Check if driver is connected
func (d *Neo4jSession) IsValid() error {
	return isValidDriver(d)
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
		return nil, errors.New(ERROR_DB_UNABLE_2_CONNECT)
	}
	
	return driver, nil
}


func isValidDriver(d *Neo4jSession) error {
	//if driver wasn't even created before
	if d == nil {
		return errors.New(ERROR_DB_MISSING_CONNECTION)
	}
	
	//if driver is set but with errors
	err := d.driver.VerifyConnectivity()
	if err != nil {
		return errors.New(ERROR_DB_MISSING_CONNECTION)
	}

	return nil
}
