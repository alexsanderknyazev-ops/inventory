package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	defaultContentTypeHeader     = "Content-Type"
	defaultApplicationTypeHeader = "application/json"

	rarityRoute = "rarity"
	nameRoute   = "name"
	idRoute     = "id"

	errorLogParceInt = "Error - parse int"
)

func getParceId(r *http.Request, logs string) (int64, error) {
	idStr := chi.URLParam(r, idRoute)
	log.Println(logs, idStr)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
