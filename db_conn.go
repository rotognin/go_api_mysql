package main

import (
	"net/http"
	"database/sql"
	"time"
	"os"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func getEnvVar(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		// var erro string
		// erro = err.Error()
		// c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Arquivo .env não carregado", "error": erro})
		return ""
	}

	return os.Getenv(key)
}

func connectDB(c *gin.Context) {
	var DB_DRIVER string = getEnvVar("DB_DRIVER")
	var DB_USER string = getEnvVar("DB_USER")
	var DB_PASS string = getEnvVar("DB_PASS")
	var DB_CONN string = getEnvVar("DB_CONN")
	var DB_DATABASE string = getEnvVar("DB_DATABASE")
	var conn_string string = fmt.Sprintf("%s:%s@%s/%s", DB_USER, DB_PASS, DB_CONN, DB_DATABASE)

	db, err := sql.Open(DB_DRIVER, conn_string)
	if err != nil {
		var erro string
		erro = err.Error()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Não foi aberto a conexão com o banco", "error": erro})
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	dbx.conn = db
}