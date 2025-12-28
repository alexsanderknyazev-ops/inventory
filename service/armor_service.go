package service

import (
	"inventory/database"
	"inventory/modules"
	"log"
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
	var weapons []modules.Weapon

	result := db.Where(whereRarity, rarity).Find(&weapons)
	log.Println("All Armor - ", len(weapons), " by rarity - ", rarity)
	return weapons, result.Error
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

	result := db.Where(whereName, name).Find(&armor)
	log.Println("GetArmorByName - Armor Name = ", armor.Name)
	return armor, result.Error
}

func CreateArmor(armor *modules.Armor) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}

	result := db.Create(armor)

	return result.Error
}
