package db

const (
	ERROR_DB_MISSING_URI         string = "cannot reach database's URI"
	ERROR_DB_CREDENTIALS         string = "user or password incorrect"
	ERROR_DB_UNABLE_2_CONNECT    string = "unable to connect to database server"
	ERROR_DB_MISSING_CONFIG_FILE string = "unable to locate .env file"
	ERROR_DB_MISSING_CONNECTION  string = "connection to database missing"
)
