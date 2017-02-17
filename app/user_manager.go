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
	usersTable := `
	CREATE TABLE IF NOT EXISTS users(
		name TEXT NOT NULL PRIMARY KEY,
		password TEXT NOT NULL,
		salt BLOB NOT NULL
	);
	`

	_, err := um.database.Exec(usersTable)
	if err != nil {
		panic(err)
	}
}

func (um *userManager) addUser(u User) {}

func (um *userManager) readUser(name string) (User, bool) {
	return defaultUser{}, false
}

func (um *userManager) updateUser(u User) {}

func (um *userManager) removeUser(name string) {}
