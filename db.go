package storage


import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


var db *sql.DB

type User struct {
	username string;
	password string
}


func GetParamsFromDB(username, password string) (User, error) {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")

	if err != nil {
		return User{"", ""}, err
	}

	var user User
	sql := "SELECT * FROM users WHERE username=?"
	err = db.QueryRow(sql, username).Scan(&user.username, &user.password)

	if err != nil {
		return User{"", ""}, err
	}

	return user, nil
}


func CreateUser(username, password string) error {
	var err error
	user, err = GetParamsFromDB(username, password)
	
	if err != nil { return err }

	if user != User{"", ""} { return error.Error("User exists") }

	return nil	
}


func createDB () error {
	var err error
	db, err = sql.Open("sqlite3", "./test.db")
	if err != nil { return err }

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

