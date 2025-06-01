package storage

import (
	"database/sql"
	"emess/user"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func GetParamsFromDB(username, password string) (user.User, error) {
	var err error
	db, err = sql.Open("sqlite3", "./db/users.db")

	if err != nil {
		return user.User{}, err
	}

	var u user.User
	s := "SELECT * FROM users WHERE username=?"
	err = db.QueryRow(s, username).Scan(&u.Username, &u.Password)

	if u.Password != password || err != nil {
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

	db, err = sql.Open("sqlite3", "./db/users.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users(username, password) VALUES (?, ?)",
		username,
		password,
	)

	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	return nil
}

func CreateDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./db/users.db")
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
