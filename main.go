package main

import (
	"fmt"
	"net/http"

	"github.com/matheuswww/mystream/src/router"
	"github.com/matheuswww/mystream/src/routes"
)

func main() {
	fmt.Println("App Running!!!")
	r := &router.Router{}
	routes.InitRoutes(r, nil)
	http.ListenAndServe(":8080", r)
}