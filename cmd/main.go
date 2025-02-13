package main

import (
	"AhmadAbdelrazik/arbun/internal/handlers"
	"log"
)

func main() {
	app := handlers.NewApplication()

	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
