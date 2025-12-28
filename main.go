package main

import (
	"log"
	"market/database"
	"market/router"
	"net/http"
)

func main() {
	database.InitDB()
	log.Default().Println("main - Init DB")

	router := router.Route()
	log.Default().Println("main - Init Route")

	port := ":8081"
	log.Printf("Server starting on %s", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}

}
