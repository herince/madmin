package app

import "database/sql"

type UserManager interface {
	CreateUser(User)
	ReadUserById(string) (User, bool)
	ReadUserByName(string) (User, bool)
	UpdateUser(User)
	RemoveUser(string)

	ValidateUser(string, string) bool
}

func NewUserManager(db *sql.DB) UserManager {
	um := &defaultUserManager{database: db}

	um.initUsersTable()

	return um
}

type defaultUserManager struct {
	database *sql.DB
}

func (um *defaultUserManager) initUsersTable() {
	usersTable := `
	CREATE TABLE IF NOT EXISTS
		users (
			id
				TEXT NOT NULL PRIMARY KEY,
			name
				TEXT NOT NULL UNIQUE,
			password
				TEXT NOT NULL,
			salt
				BLOB NOT NULL
	);
	`

	_, err := um.database.Exec(usersTable)
	if err != nil {
		panic(err)
	}
}

func (um *defaultUserManager) CreateUser(u User) {
	stmt, err := um.database.Prepare(`
		INSERT INTO
			users (
				id,
				name,
				password,
				salt)
		VALUES(?, ?, ?, ?)
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		u.ID(),
		u.Name(),
		u.Password(),
		u.Salt())
	if err != nil {
		panic(err)
	}
}

func (um *defaultUserManager) ReadUserById(id string) (User, bool) {
	stmt, err := um.database.Prepare(`
	SELECT
		name,
		password,
		salt
	FROM
		users
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	u := defaultUser{id: id}
	err = stmt.QueryRow(id).Scan(
		&u.name,
		&u.password,
		&u.salt)
	switch {
	case err == sql.ErrNoRows:
		return nil, false
	case err != nil:
		panic(err)
	}

	return &u, true
}

func (um *defaultUserManager) ReadUserByName(name string) (User, bool) {
	stmt, err := um.database.Prepare(`
	SELECT
		id,
		password,
		salt
	FROM
		users
	WHERE
		name = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	u := defaultUser{name: name}
	err = stmt.QueryRow(name).Scan(
		&u.id,
		&u.password,
		&u.salt)
	switch {
	case err == sql.ErrNoRows:
		return nil, false
	case err != nil:
		panic(err)
	}

	return &u, true
}

func (um *defaultUserManager) UpdateUser(u User) {
	stmt, err := um.database.Prepare(`
	UPDATE
		users
	SET
		name = ?,
		password = ?,
		salt = ?
	WHERE
		id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(
		u.Name(),
		u.Password(),
		u.Salt(),
		u.ID())
	if err != nil {
		panic(err)
	}
}

func (um *defaultUserManager) RemoveUser(id string) {
	stmt, err := um.database.Prepare(`
		DELETE FROM
			users
		WHERE
			id = ?
	`)
	defer stmt.Close()
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		panic(err)
	}
}

func (um *defaultUserManager) ValidateUser(name, password string) bool {
	u, ok := um.ReadUserByName(name)
	if ok {
		return u.CheckPassword(password)
	}
	return false
}
