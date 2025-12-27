package router

import (
	"inventory/handler"

	"github.com/go-chi/chi/v5"
)

func Route() *chi.Mux {

	r := chi.NewRouter()

	r.Route("/inventory", func(r chi.Router) {
		r.Get("/weapons", handler.GetAllWeapon)
		r.Get("/armor", handler.GetAllArmor)
		r.Post("/weapons", handler.CreateWeapon)
		r.Post("/armor", handler.CreateArmor)
		// r.Post("/", handler.CreateUser)
		// r.Get("/{id}", handler.GetUserById)
		// r.Get("/limit/{limit}", handler.GetTopUsersByLimit)
		// // r.Put("/{id}", handler.UpdateClient)
		// // r.Delete("/{id}", handler.DeleteClient)
	})
	return r
}
