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

func GetAllWeaponByRarity(rarity string) ([]modules.Weapon, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}
	var weapons []modules.Weapon

	result := db.Where(whereRarity, rarity).Find(&weapons)
	log.Println("AllWeapon -", len(weapons), " by rarity ", rarity)
	return weapons, result.Error
}

func GetWeaponById(id int64) (modules.Weapon, error) {
	db := database.GetDB()
	var weapon modules.Weapon
	result := db.First(&weapon, id)
	log.Println("GetWeaponById - Weapon ID = ", weapon.ID)
	return weapon, result.Error
}

func GetWeaponByName(name string) (modules.Weapon, error) {
	db := database.GetDB()
	var weapon modules.Weapon
	result := db.Where(whereName, name).Find(&weapon)
	log.Println("GetWeaponByName - Weapon Name = ", weapon.Name)
	return weapon, result.Error
}

func CreateWeapon(weapon *modules.Weapon) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}

	result := db.Create(weapon)

	return result.Error
}
