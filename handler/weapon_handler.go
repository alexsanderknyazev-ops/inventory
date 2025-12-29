package handler

import (
	"encoding/json"
	"log"
	"market/modules"
	"market/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllWeapon(w http.ResponseWriter, r *http.Request) {
	result, err := service.GetAllWeapon()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func GetAllWeaponByRarity(w http.ResponseWriter, r *http.Request) {
	rarity := chi.URLParam(r, rarityRoute)
	result, err := service.GetAllWeaponByRarity(rarity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func GetWeaponByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, nameRoute)
	log.Println("GetWeaponByName - name = ", name)
	result, err := service.GetWeaponByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if result.ID == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func GetWeaponById(w http.ResponseWriter, r *http.Request) {

	id, err := getParceId(r, "GetWeaponById - idStr = ")
	if err != nil {
		log.Println(errorLogParceInt)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := service.GetWeaponById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func CreateWeapon(w http.ResponseWriter, r *http.Request) {
	var newWeapon modules.Weapon
	if err := json.NewDecoder(r.Body).Decode(&newWeapon); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := service.CreateWeapon(&newWeapon)

	if result != nil {
		http.Error(w, result.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWeapon)

}

func CreateWeaponBatch(w http.ResponseWriter, r *http.Request) {
	var newWeapons []modules.Weapon

	if err := json.NewDecoder(r.Body).Decode(&newWeapons); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newWeapons) == 0 {
		http.Error(w, "No armor data provided", http.StatusBadRequest)
		return
	}

	log.Printf("Creating %d armors in batch", len(newWeapons))
	_, err := service.CreateWeaponBatch(newWeapons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Armors created successfully",
		"count":   len(newWeapons),
		"armors":  newWeapons,
	})
}

func DeleteWeaponById(w http.ResponseWriter, r *http.Request) {
	id, err := getParceId(r, "DeleteWeaponById - idStr = ")
	if err != nil {
		log.Println(errorLogParceInt)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.DeleteWeaponById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteWeaponByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, nameRoute)
	err := service.DeleteWeaponByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
