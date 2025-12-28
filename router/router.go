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

	postWeapon = "/weapons"
	postArmor  = "/armor"
)
