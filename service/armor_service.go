package service

import (
	"log"
	"market/database"
	"market/modules"

	"github.com/IBM/sarama"

	"gorm.io/gorm"
)

func GetAllArmor() ([]modules.Armor, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}
	var armor []modules.Armor

	result := db.Find(&armor)

	return armor, result.Error
}

func GetAllArmorByRarity(rarity string) ([]modules.Weapon, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}
	var armors []modules.Weapon

	result := db.Where(whereRarity, rarity).Find(&armors)
	log.Println("All Armor - ", len(armors), " by rarity - ", rarity)
	return armors, result.Error
}

func GetArmorById(id int64) (modules.Armor, error) {
	db := database.GetDB()

	var armor modules.Armor

	result := db.First(&armor, id)
	log.Println("GetArmorById - Armor ID = ", armor.ID)
	return armor, result.Error
}

func GetArmorByName(name string) (modules.Armor, error) {
	db := database.GetDB()
	var armor modules.Armor

	result := db.Where("name = ?", name).First(&armor)
	if result.Error != nil {
		log.Printf("GetArmorByName - Error finding armor '%s': %v", name, result.Error)
		return armor, result.Error
	}

	log.Println("GetArmorByName - Found Armor Name = ", armor.Name)

	if armor.Name != "" {
		sendToKafkaDirectly(armor.Name)
	}

	return armor, nil
}

func sendToKafkaDirectly(message string) error {
	log.Printf("Sending to Kafka: %s", message)

	brokers := []string{"kafka-service:9092"}
	topic := "market-events"

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Printf("Failed to create Kafka producer: %v", err)
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Kafka send failed: %v", err)
		return err
	}

	log.Printf("Kafka send successful: partition=%d, offset=%d", partition, offset)
	return nil
}

func GetArmorByPrice(price float64, is bool) ([]modules.Armor, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}
	var armors []modules.Armor
	var res *gorm.DB
	if is {
		res = db.Where("price > ?", price).Find(&armors)
	} else {
		res = db.Where("price < ?", price).Find(&armors)
	}

	result := res.
		Order("price DESC").
		Find(&armors)

	return armors, result.Error
}

func CreateArmor(armor *modules.Armor) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}

	result := db.Create(armor)

	return result.Error
}

func CreateArmorBatch(armors []modules.Armor) ([]modules.Armor, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}

	if len(armors) == 0 {
		return nil, nil
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.CreateInBatches(&armors, 100)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return armors, nil
}

func DeleteArmorById(id int64) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}
	result := db.Delete(&modules.Armor{}, id)
	return result.Error
}

func DeleteArmorByName(name string) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}

	result := db.Where(whereName, name).Delete(&modules.Armor{})

	return result.Error
}
