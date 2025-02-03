package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/matheuswww/mystream/src/logger"
)

func NewSql() *sql.DB {
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PASSWORD")
	name := os.Getenv("SQL_DB_NAME")
	port := os.Getenv("SQL_PORT")
	host := os.Getenv("HOST")
	if user == "" || password == "" || name == "" || port == "" || host == "" {
		log.Fatalf("Error trying get env, user: %s, password: %s, name: %s, port: %s", user, password, name, port)
	}
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, name, host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error trying conn to sql: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("error trying conn to sql: %v", err)
	}
	logger.Log("Sql is running!!!")
	return db
}
