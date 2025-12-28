package router

import (
	"inventory/handler"

	"github.com/go-chi/chi/v5"
)

const (
	inventoryRoute       = "/inventory"
	getAllWeapon         = "/weapons"
	getAllArmor          = "/armors"
	getAllWeaponByRarity = "/weapons/rarity/{rarity}"
	getAllArmorByRarity  = "/armors/rarity/{rarity}"
	getArmorByName       = "/armors/name/{name}"
	getWeaponsByName     = "/weapons/name/{name}"
	getArmorById         = "/armors/{id}"
	getWeaponsById       = "/weapons/{id}"

	postWeapon     = "/weapons"
	postArmor      = "/armors"
	postBatchArmor = "/armors/batch"

	deleteWeaponByName = "/weapons/name/{name}"
	deleteArmorByName  = "/armors/name/{name}"
	deleteWeaponById   = "/weapons/{id}"
	deleteArmorById    = "/armors/{id}"
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
		r.Post(postBatchArmor, handler.CreateArmorBatch)

		r.Delete(deleteWeaponByName, handler.DeleteWeaponByName)
		r.Delete(deleteArmorByName, handler.DeleteArmorByName)
		r.Delete(deleteWeaponById, handler.DeleteWeaponById)
		r.Delete(deleteArmorById, handler.DeleteArmorById)
	})
	return r
}
