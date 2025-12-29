package handler

import (
	"encoding/json"
	"log"
	"market/modules"
	"market/service"
	"net/http"
	"strconv"

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
func GetArmorByPrice(w http.ResponseWriter, r *http.Request) {

	priceStr := chi.URLParam(r, "price")
	log.Printf("Raw price string: '%s'", priceStr)
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("ParseFloat error: %v", err)
		log.Printf("Trying to parse '%s'", priceStr)
	}
	isStr := chi.URLParam(r, boolPrice)
	is, err := strconv.ParseBool(isStr)
	log.Println(is)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := service.GetArmorByPrice(price, is)
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

func CreateArmorBatch(w http.ResponseWriter, r *http.Request) {
	var newArmors []modules.Armor
	if err := json.NewDecoder(r.Body).Decode(&newArmors); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newArmors) == 0 {
		http.Error(w, "No armor data provided", http.StatusBadRequest)
		return
	}

	log.Printf("Creating %d armors in batch", len(newArmors))
	_, err := service.CreateArmorBatch(newArmors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Armors created successfully",
		"count":   len(newArmors),
		"armors":  newArmors,
	})
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
