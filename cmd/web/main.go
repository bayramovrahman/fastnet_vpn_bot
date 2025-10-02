package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/config"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/handlers"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/render"
)

const portNumber = ":8080"

func main() {
	var app config.AppConfig

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Starting serve on port %s", portNumber)

	serve := &http.Server{
		Addr: portNumber,
		Handler: routes(),
	}

	err = serve.ListenAndServe()
	log.Fatal(err)
}
