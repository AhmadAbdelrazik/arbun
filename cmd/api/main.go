package main

import (
	"AhmadAbdelrazik/arbun/cmd/api/app"
	"log"
)

func main() {
	app := app.NewApplication()

	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
