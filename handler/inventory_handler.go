package handler

import (
	"encoding/json"
	"inventory/modules"
	"inventory/service"
	"net/http"
	// "strconv"
	// "github.com/go-chi/chi/v5"
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

func GetAllArmor(w http.ResponseWriter, r *http.Request) {
	result, err := service.GetAllWeapon()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func CreateArmor(w http.ResponseWriter, r *http.Request) {
	var newArmor modules.Armor
	if err := json.NewDecoder(r.Body).Decode(&newArmor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result := service.CreateArmor(&newArmor)

	if result != nil {
		http.Error(w, result.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newArmor)

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

const (
	defaultContentTypeHeader     = "Content-Type"
	defaultApplicationTypeHeader = "application/json"
)
