package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	db "github.com/matheuswww/mystream/src/configuration/sql"
	"github.com/matheuswww/mystream/src/logger"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
	db := db.NewSql()
	id := uuid.NewString()
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")
	name :=  os.Getenv("ADMIN_NAME")
	if email == "" || password == "" || name == "" {
		log.Fatal("Error trying get env")
	}
	encryptedPassword,err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Fatal(err)
	}
	query := "INSERT INTO admin (id, email, name, password) VALUES ($1, $2, $3, $4)"
	_,err = db.ExecContext(ctx, query, id, email, name, encryptedPassword)
	if err != nil {
		log.Fatal(err)
	}
	logger.Log("User inserted with success")
}