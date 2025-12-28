package router

import (
	"inventory/handler"

	"github.com/go-chi/chi/v5"
)

func Route() *chi.Mux {

	r := chi.NewRouter()

	r.Route(inventoryRoute, func(r chi.Router) {
		r.Get(getAllWeapon, handler.GetAllWeapon)
		r.Get(getAllArmor, handler.GetAllArmor)
		r.Get(getAllWeaponByRarity, handler.GetAllWeaponByRarity)
		r.Get(getAllArmorByRarity, handler.GetAllArmorByRarity)
		r.Get(getArmorByName, handler.GetArmorByName)
		r.Get(getWeaponsByName, handler.GetWeaponByName)
		r.Get(getArmorById, handler.GetArmorById)
		r.Get(getWeaponsById, handler.GetWeaponById)

		r.Post(postWeapon, handler.CreateWeapon)
		r.Post(postArmor, handler.CreateArmor)

	})
	return r
}

const (
	inventoryRoute = "/inventory"
	getAllWeapon   = "/weapons"
	getAllArmor    = "/armors"

	getAllWeaponByRarity = "/weapons/rarity/{rarity}"
	getAllArmorByRarity  = "/armors/rarity/{rarity}"

	getArmorByName   = "/armors/name/{name}"
	getWeaponsByName = "/weapons/name/{name}"

	getArmorById   = "/armors/{id}"
	getWeaponsById = "/weapons/{id}"

	postWeapon = "/weapons"
	postArmor  = "/armor"
)
