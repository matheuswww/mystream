package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/matheuswww/mystream/src/router"
	"github.com/matheuswww/mystream/src/routes"
)

func main() {
	fmt.Println("App Running!!!")
	initEnv()
	r := &router.Router{}
	r.Middleware(middleware)
	routes.InitRoutes(r, nil)
	http.ListenAndServe(":8080", r)
}

func initEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func middleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")	
}