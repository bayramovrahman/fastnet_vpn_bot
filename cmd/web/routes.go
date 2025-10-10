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
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Login)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	
	mux.Group(func(r chi.Router) {
		r.Use(Auth)
		r.Get("/home", handlers.Repo.Home)
		r.Get("/taxes", handlers.Repo.Taxes)
		r.Get("/logout", handlers.Repo.Logout)
		r.Get("/invoice", handlers.Repo.Invoice)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
