package main

import (
	"log"

	"github.com/herince/madmin/app"
)

func main() {
	var (
		port = ":4200"
		dbPath = "./database/database.sqlite"
	)

	server := app.NewMAdminServer(port, dbPath)

	log.Println("Listening...")
	err := server.ListenAndServe();
	if err != nil {
		panic(err)
	}
}
