package main

import (
	"AhmadAbdelrazik/arbun/internal/api"
	"log"
)

func main() {
	app := handlers.NewApplication()

	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
