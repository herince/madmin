package main

import (
	"os"
	"database/sql"

	"github.com/herince/madmin/app"
)

func main () {
	db, err := sql.Open("sqlite3", "./database/database.sqlite")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	user, err := app.NewUser(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}

	um := app.NewUserManager(db)
	um.CreateUser(user)
	
	someUser, exists := um.ReadUserByName(os.Args[1])
	if !exists {
		print("Oh, noes!")
	} else {
		print(someUser.Name(), " lives!\n")
	}
}