package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	db "github.com/matheuswww/mystream/src/configuration/sql"
	"github.com/matheuswww/mystream/src/routes"
)

func main() {
	fmt.Println("App Running!!!")
	initEnv()
	db := db.NewSql()
	r := gin.Default()
	r.Use(cors.Default())
	routes.InitRoutes(r, db)
	http.ListenAndServe(":8080", r)
}

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func configCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://127.0.0.1:5000"},
    AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
	}))
}