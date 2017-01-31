package main

import (
	"github.com/herince/madmin/app"
	"log"
	"net/http"
)

func main() {
	// 	fs := http.FileServer(http.Dir("static/"))
	// 	http.Handle("/", fs)

	mah := app.NewMAdminHandler()
	http.Handle("/", mah)

	log.Println("Listening...")
	http.ListenAndServe(":4200", nil)
}
