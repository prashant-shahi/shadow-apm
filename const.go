package main

import "os"

// Config variables
const (
	apmServerIPAddress = "http://localhost"
	apmServerPort      = "8200"
)

// DB variable
const (
	COLLECTION = "transactions"
	DATABASE = "apm"
	DB_SERVER_PORT = os.Getenv("DB_SERVER_PORT")
	DB_SERVER_URL = os.Getenv("DB_SERVER_URL")
	DBUSER = os.Getenv("DBUSER")
	DBPWD = os.Getenv("DBPWD")
)