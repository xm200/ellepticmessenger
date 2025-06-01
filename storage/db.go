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
	DB, err = sql.Open("sqlite3", "./db/users.db")

	if err != nil {
		return user.User{}, err
	}

	var u user.User
	s := "SELECT * FROM users WHERE username=?"
	err = DB.QueryRow(s, username).Scan(&u.Username, &u.Password)

	if u.Password != password || err != nil {
		return user.User{}, err
	}
	defer DB.Close()

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

	_, err = tx.Exec("INSERT INTO users(username, password) VALUES (?, ?)", u.Username, u.Password)
	tx.Commit()
	if err != nil {
		log.Println(err)
	}

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
 	        username TEXT NOT NULL,
        	password TEXT NOT NULL
    	);`

	_, err = tx.Exec(create)
	tx.Commit()
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
	return nil
}
