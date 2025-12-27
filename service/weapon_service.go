package service

import (
	"inventory/database"
	"inventory/modules"
	"log"
)

func GetAllWeapon() ([]modules.Weapon, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}
	var weapons []modules.Weapon

	result := db.Find(&weapons)
	log.Println("AllWeapon -", len(weapons))
	return weapons, result.Error
}

func CreateWeapon(weapon *modules.Weapon) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}

	result := db.Create(weapon)

	return result.Error
}
