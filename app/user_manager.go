package app

import "database/sql"

type userManager struct {
	database *sql.DB
}

func newUserManager(db *sql.DB) *userManager {
	um := &userManager{database: db}

	um.initUsersTable()

	return um
}

func (um *userManager) initUsersTable() {
	users_table := `
	CREATE TABLE IF NOT EXISTS users(
		id BLOB NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		password TEXT NOT NULL
	);
	`

	_, err := um.database.Exec(users_table)
	if err != nil {
		panic(err)
	}
}
