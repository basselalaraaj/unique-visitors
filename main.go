package main

import (
	"github.com/basselalaraaj/unique-visitors/routes"
)

func main() {
	router := routes.GetRoutes()
	router.Run(":8080")
}
