package storage

import (
	"database/sql"
	"emess/user"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

func GetParamsFromDB(username, password string) (user.User, error) {
	var err error
	log.Println(username, password)
	DB, err = sql.Open("sqlite3", "./db/users.db")

	if err != nil {
		return user.User{}, err
	}

	var u user.User
	err = DB.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&u.Username, &u.Password)
	defer DB.Close()

	if u.Password != password || err != nil {
		log.Println(err)
		return user.User{}, err
	}

	return u, nil
}

func CreateUser(username, password string) {
	u, err := GetParamsFromDB(username, password)

	if err != nil {
		log.Println(err)
	}

	if user.Equal(u, user.User{}) {
		log.Println(err)
	}

	DB, err = sql.Open("sqlite3", "./db/users.db")
	if err != nil {
		log.Println(err)
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
	}

	_, err = tx.Exec("INSERT INTO users(username, password) VALUES (?, ?)", username, password)
	tx.Commit()

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(DB)
}

func CreateDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./db/users.db")
	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
	}

	create := `
	CREATE TABLE users (
 	        username VARCHAR NOT NULL,
        	password VARCHAR NOT NULL
    	);`

	_, err = tx.Exec(create)
	tx.Commit()
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
	return nil
}
