package main

import (
	"log"
	"market/database"
	"market/router"
	"net/http"
	"time"
)

func main() {

	database.InitDB()
	log.Println("main - Init DB")

	router := router.Route()
	log.Println("main - Init Route")

	port := "8070"
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on :%s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
