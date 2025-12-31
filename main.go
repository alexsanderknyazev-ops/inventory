package main

import (
	"log"
	"market/config"
	"market/database"
	"market/kafka"
	"market/router"
	"net/http"
	"time"
)

func main() {

	config.InitKafka()
	log.Println("Kafka config loaded")

	// 2. Инициализация продюсера
	if err := kafka.Init(); err != nil {
		log.Printf("Warning: Kafka producer not connected: %v", err)
		log.Println("Running without Kafka...")
	} else {
		defer kafka.Global.Close()
		log.Println("Kafka producer ready")
	}

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
