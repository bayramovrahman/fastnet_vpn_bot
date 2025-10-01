package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	r := routes()

	fmt.Println("Server running at http://localhost" + portNumber)
	err := http.ListenAndServe(portNumber, r)
	if err != nil {
		log.Fatal(err)
	}
}
