package kafka

import (
	"encoding/json"
	"log"
	"market/config"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

var Global *Producer

func Init() error {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, cfg)
	if err != nil {
		return err
	}

	Global = &Producer{
		syncProducer: producer,
		topic:        config.Kafka.Topic,
	}

	log.Printf("Kafka producer connected to brokers: %v", config.Kafka.Brokers)
	return nil
}

func (p *Producer) Send(key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(jsonData),
		Timestamp: time.Now(),
	}

	_, _, err = p.syncProducer.SendMessage(msg)
	return err
}

func (p *Producer) SendString(message string) error {
	log.Printf("Kafka: Sending to topic '%s': %s", p.topic, message)

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)

	if err != nil {
		log.Printf("Kafka send failed: %v", err)
	} else {
		log.Printf("Kafka send successful: partition=%d, offset=%d", partition, offset)
	}

	return err
}

func (p *Producer) Close() error {
	if p.syncProducer != nil {
		return p.syncProducer.Close()
	}
	return nil
}
