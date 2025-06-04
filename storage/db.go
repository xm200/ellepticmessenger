package storage

import (
	"database/sql"
	"emess/user"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
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
	err = DB.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&u.Id, &u.Username, &u.Password)
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

type Nickname struct {
	Nick string
}

func OnlineAdd(username string) {
	var err error
	DB, err = sql.Open("sqlite3", "./db/online.db")
	if err != nil {
		log.Println(err)
		return
	}

	r, err := DB.Exec("INSERT INTO online(username) VALUES(?)", username)
	defer DB.Close()

	if err != nil {
		log.Println(err)
		return
	}
	log.Println(r.RowsAffected())
}

func OnlineDelete(username string) {
	time.Sleep(5 * 60 * time.Second)
	var err error
	DB, err = sql.Open("sqlite3", "./db/online.db")
	if err != nil {
		log.Println(err)
		return
	}
	r, _ := DB.Exec("DELETE FROM online WHERE username = ?", username)
	log.Println(r.RowsAffected())
	defer DB.Close()
}

func OnlineGet() []Nickname {
	var err error
	DB, err = sql.Open("sqlite3", "./db/online.db")
	if err != nil {
		log.Println(err)
		return nil
	}
	var users []Nickname
	rows, err := DB.Query("SELECT username FROM online")

	if err != nil {
		log.Println(err)
		return nil
	}

	for rows.Next() {
		var username string
		rows.Scan(&username)
		users = append(users, Nickname{username})
	}

	defer DB.Close()
	return users
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
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
 	        username VARCHAR NOT NULL UNIQUE ,
        	password VARCHAR NOT NULL
    	);`

	_, err = tx.Exec(create)
	tx.Commit()
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
	DB, err = sql.Open("sqlite3", "./db/online.db")
	tx, err = DB.Begin()
	if err != nil {
		log.Println(err)
	}

	create = `
	CREATE TABLE online (
	    	id INTEGER PRIMARY KEY AUTOINCREMENT,
 	        username VARCHAR NOT NULL UNIQUE
    	);`

	_, err = tx.Exec(create)
	tx.Commit()
	if err != nil {
		log.Println(err)
	}
	defer DB.Close()
	return nil
}
