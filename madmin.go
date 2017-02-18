package main

import (
	"github.com/herince/madmin/app"
)

func main() {
	var (
		port = ":4200"
		dbPath = "./database/database.sqlite"
	)
	app.Init(port, dbPath)
}
