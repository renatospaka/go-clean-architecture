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
