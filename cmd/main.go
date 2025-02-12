package main

import (
	"AhmadAbdelrazik/arbun/internal/handlers"
	"log"
)

func main() {
	app := handlers.New()

	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
