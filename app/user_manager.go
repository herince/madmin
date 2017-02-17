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
		name TEXT NOT NULL PRIMARY KEY,
		password TEXT NOT NULL
	);
	`

	_, err := um.database.Exec(users_table)
	if err != nil {
		panic(err)
	}
}

func (um *userManager) addUser(u user) {}

func (um *userManager) readUser(name string) (user, bool) {
	return user{}, false
}

func (um *userManager) updateUser(u user) {}

func (um *userManager) removeUser(name string) {}
