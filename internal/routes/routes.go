package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/nhx-finance/kesy/internal/app"
)


func SetUpRoutes(kesy *app.Application) *chi.Mux {
	r := chi.NewRouter()

	/**
	* APPLICATION STATUS ROUTES
	*/
	r.Get("/status", kesy.HandleStatus)



	return r
}