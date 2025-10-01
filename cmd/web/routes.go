package main

import (
	"net/http"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	// "/" veya "/home" açıldığında Home handler çalışacak
	mux.Get("/", handler.Home)
	mux.Get("/home", handler.Home)

	return mux
}
