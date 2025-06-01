package storage

import (
	"database/sql"
	"emess/user"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func GetParamsFromDB(username, password string) (user.User, error) {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")

	if err != nil {
		return user.User{}, err
	}

	var u user.User
	s := "SELECT * FROM users WHERE username=?"
	err = db.QueryRow(s, username).Scan(&u.Username, &u.Password)

	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func CreateUser(username, password string) error {
	u, err := GetParamsFromDB(username, password)

	if err != nil {
		return err
	}

	if user.Equal(u, user.User{}) {
		return err
	}

	return nil
}

func createDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		return err
	}

	create := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
 	        username TEXT NOT NULL UNIQUE,
        	password TEXT NOT NULL
    	);`

	_, err = db.Exec(create)
	if err != nil {
		return err
	}

	return nil
}
