package main

import (
	"net/http"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/home", handlers.Repo.Home)
	mux.Get("/taxes", handlers.Repo.Taxes)
	mux.Get("/invoice", handlers.Repo.Invoice)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
