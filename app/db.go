package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func main () {
	db, err := sql.Open("postgres", "user=madmin password=roni dbname=madmin sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	name := "roni"
	rows, err := db.Query("SELECT id, salary FROM customers WHERE name = $1", name)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var salaryString string
		err = rows.Scan(&id, &salaryString)
		salary, err := decimal.NewFromString(salaryString)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d gets %s\n", id, salary)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
}
