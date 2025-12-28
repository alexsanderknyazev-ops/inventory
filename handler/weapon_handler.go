package handler

import (
	"encoding/json"
	"inventory/modules"
	"inventory/service"
	"log"
	"net/http"
	"strconv"

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
	idStr := chi.URLParam(r, idRoute)
	log.Println("GetWeaponById - idStr = ", idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
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
