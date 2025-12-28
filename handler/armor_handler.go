package handler

import (
	"encoding/json"
	"inventory/modules"
	"inventory/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllArmor(w http.ResponseWriter, r *http.Request) {
	result, err := service.GetAllArmor()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func GetAllArmorByRarity(w http.ResponseWriter, r *http.Request) {
	rarity := chi.URLParam(r, rarityRoute)
	result, err := service.GetAllArmorByRarity(rarity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(defaultContentTypeHeader, defaultApplicationTypeHeader)
	json.NewEncoder(w).Encode(result)
}

func GetArmorByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, nameRoute)
	log.Println("GetArmorByName - name = ", name)
	result, err := service.GetArmorByName(name)
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

func GetArmorById(w http.ResponseWriter, r *http.Request) {
	id, err := getParceId(r, "GetArmorById - idStr = ")
	if err != nil {
		log.Panicln(errorLogParceInt)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := service.GetArmorById(id)
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

func DeleteArmorById(w http.ResponseWriter, r *http.Request) {
	id, err := getParceId(r, "DeleteArmorById - idStr = ")
	if err != nil {
		log.Println(errorLogParceInt)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.DeleteArmorById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteArmorByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, nameRoute)

	err := service.DeleteArmorByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
