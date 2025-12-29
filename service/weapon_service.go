package service

import (
	"log"
	"market/database"
	"market/modules"
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

func CreateWeaponBatch(weapons []modules.Weapon) ([]modules.Weapon, error) {
	db := database.GetDB()
	if db == nil {
		return nil, nil
	}

	if len(weapons) == 0 {
		return nil, nil
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.CreateInBatches(&weapons, 100)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return weapons, nil
}

func DeleteWeaponById(id int64) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}
	result := db.Delete(&modules.Weapon{}, id)
	return result.Error
}

func DeleteWeaponByName(name string) error {
	db := database.GetDB()
	if db == nil {
		return nil
	}
	result := db.Where(whereName, name).Delete(&modules.Weapon{})

	return result.Error
}
