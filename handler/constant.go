package handler

import (
	"log"
	"strconv"
)

const (
	defaultContentTypeHeader     = "Content-Type"
	defaultApplicationTypeHeader = "application/json"

	rarityRoute = "rarity"
	nameRoute   = "name"
	idRoute     = "id"

	errorLogParceInt = "Error - parse int"
)

func getParceId(str string, logs string) (int64, error) {
	log.Println(logs, str)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
