package main

import (
	"log"
	"market/config"
	"market/database"
	"market/grpc"
	"market/kafka"
	"market/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Канал для graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	config.InitKafka()
	log.Println("Kafka config loaded")

	// Инициализация продюсера
	if err := kafka.Init(); err != nil {
		log.Printf("Warning: Kafka producer not connected: %v", err)
		log.Println("Running without Kafka...")
	} else {
		defer kafka.Global.Close()
		log.Println("Kafka producer ready")
	}

	database.InitDB()
	log.Println("main - Init DB")

	// Запускаем GRPC сервер в отдельной горутине
	grpcErr := make(chan error, 1)
	go func() {
		port := "50061" // GRPC порт для market
		log.Printf("Starting Market GRPC server on :%s", port)
		grpcErr <- grpc.StartServer(port)
	}()

	// Даем время на запуск GRPC сервера
	time.Sleep(100 * time.Millisecond)

	// Запускаем HTTP сервер
	httpErr := make(chan error, 1)
	go func() {
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

		log.Printf("HTTP Server starting on :%s", port)
		httpErr <- server.ListenAndServe()
	}()

	// Ждем сигнал завершения или ошибку
	select {
	case sig := <-stopChan:
		log.Printf("Received signal: %v. Shutting down...", sig)

		// Graceful shutdown (просто ждем)
		time.Sleep(2 * time.Second)
		log.Println("Graceful shutdown complete")

	case err := <-httpErr:
		log.Printf("HTTP server error: %v", err)
	case err := <-grpcErr:
		log.Printf("GRPC server error: %v", err)
	}

	log.Println("Market service shutdown complete")
}
