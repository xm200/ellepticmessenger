package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"db.go"
	"httpserver/db.go"

	_ "github.com/mattn/go-sqlite3"
)


func generateKeyPair(w http.ResponseWriter, r *http.Request) {
	serverCurve := ecdh.X25519()
	ClientPrivKey, err := serverCurve.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(w, "Error occured, try later")
		return
	}

	PrivKey := ""
	for _, b := range ClientPrivKey.Bytes() {
		PrivKey = PrivKey + fmt.Sprintf("%x", b)
	}

	PubKey := ""
	for _, b := range ClientPrivKey.PublicKey().Bytes() {
		PubKey = PubKey + fmt.Sprintf("%x", b)
	}

	fmt.Fprintf(w, "{\"PrivKey\":\"%v\",\"Pubkey\":\"%v\"}", PrivKey, PubKey)
	log.Println("Generated keypair")
}


func Login(w http.ResponseWriter, r *http.Request) {
	if (r.FormValue("username") == "" && r.FormValue("password") == "") || (r.FormValue("email") == "" && r.FormValue("password") == "") {
		fmt.Fprintf(w, "Provide password, and email or username")
		return
	}

	password := r.FormValue("password")
	username := ""

	if r.FormValue("email") == "" {
		username = r.FormValue("username") 	
	} else {
		username = r.FormValue("email")
	}

	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Println("Fatal db error")
		return
	}
	
	rows, err := 

	if err != nil {
		log.Println("SQL Error while logging in")
	}

	if rows. [0] != password {
		http.Error(w, "There is no such user or no such password")
	}

	db.Close()
}


func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("email") == "" || r.URL.Query().Get("username") == "" || r.URL.Query().Get("password") == "" {
		fmt.Fprintf(w, "Provide email, username, password")
		return
	}
	
	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Println("Fatal db error")
		return
	}
	
	_, err = db.Exec("insert into users (email, username, password) values ($3, $1, $2)", r.URL.Query().Get("email"), r.URL.Query().Get("username"), r.URL.Query().Get("password"))
	if err != nil {
		log.Println("Fatal insert error")
	}

	err = db.Close()
	if err != nil {
		log.Println("Fatal insert error")
	}

}


func DisplayLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Will be right back")
}


func DisplayRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Will be right back")
}


func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodGet: { DisplayLogin(w, r) }
			case http.MethodPost: { Login(w, r) }
			default: http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		} 
	})
	
	mux.HandleFunc("/register", func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case http.MethodGet: { DisplayRegister(w, r) }
			case http.MethodPost: { Register(w, r) }
			default: http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})


	mux.HandleFunc("/keypair", func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			generateKeyPair(w, r)
		} else {
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})


	cfg := &tls.Config{
		MinVersion: tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521},
	}

	srv := &http.Server{
		Addr: ":443",
		Handler: mux,
		TLSConfig: cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Listening on :443...")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "keypair.pem"))
}

