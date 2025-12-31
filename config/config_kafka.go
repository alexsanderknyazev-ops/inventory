package config

import (
	"log"
	"os"
	"strings"
)

// KafkaConfig - конфигурация
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

var Kafka *KafkaConfig

func InitKafka() {
	Kafka = &KafkaConfig{
		Brokers: getBrokers(),
		Topic:   getTopic(),
	}

	log.Printf("=== KAFKA CONFIG ===")
	log.Printf("Brokers: %v", Kafka.Brokers)
	log.Printf("Topic: %s", Kafka.Topic)
	log.Printf("Env KAFKA_BROKERS: %s", os.Getenv("KAFKA_BROKERS"))
}

func getBrokers() []string {
	// 1. Всегда сначала проверяем переменную окружения
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers != "" {
		log.Printf("Using KAFKA_BROKERS from env: %s", brokers)
		return strings.Split(brokers, ",")
	}

	// 2. Автоматическое определение
	//    Проверяем, находимся ли мы внутри Kubernetes
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		// Внутри Kubernetes - используем service name
		log.Println("Detected Kubernetes environment, using kafka-service:9092")
		return []string{"kafka-service:9092"}
	} else {
		// Локальная разработка - используем localhost
		log.Println("Detected local environment, using localhost:9092")
		return []string{"localhost:9092"}
	}
}

func getTopic() string {
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		return "market-events"
	}
	return topic
}
