package main

import (
	"database/sql"
)

var dbx dbConnection

type dbConnection struct {
	conn *sql.DB
}