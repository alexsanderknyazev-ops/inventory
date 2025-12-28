package database

import (
	"log"

	"market/config"
	"market/modules"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbConfig := config.LoadDBConfig()

	var err error
	DB, err = config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	err = DB.AutoMigrate(&modules.Weapon{})
	err = DB.AutoMigrate(&modules.Armor{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}
