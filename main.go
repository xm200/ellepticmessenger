package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/tls"
	"emess/storage"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

func generateKeyPair() (string, string) {
	serverCurve := ecdh.X25519()
	PrivateKey, err := serverCurve.GenerateKey(rand.Reader)
	if err != nil {
		return "", ""
	}

	PrivKey := ""
	for _, b := range PrivateKey.Bytes() {
		PrivKey = PrivKey + fmt.Sprintf("%x", b)
	}

	PubKey := ""
	for _, b := range PrivateKey.PublicKey().Bytes() {
		PubKey = PubKey + fmt.Sprintf("%x", b)
	}

	log.Println("Generated keypair")
	return PrivKey, PubKey
}

func LoginFront(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/login.html")
	if err != nil {
		return
	}

	tmpl.Execute(w, nil)

	log.Println("Somebody want to login")
}

func LoginBack(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Authentication failed", http.StatusForbidden)
	}
	u, err := storage.GetParamsFromDB(username, password)

	if err != nil {
		log.Println(err)
		log.Printf("%v %v\n", u.Username, u.Password)
		http.Error(w, "Authentication failed", http.StatusForbidden)
		return
	}
	log.Printf("%v %v\n", u.Username, u.Password)
	priv, pub := generateKeyPair()
	_, err = fmt.Fprintf(w, "Your keypair is %s %s", priv, pub)
	if err != nil {
		return
	}
}

func RegisterFront(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/register.html")
	if err != nil {
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}

	log.Println("Somebody want to register")
}

func RegisterBack(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Authentication failed", http.StatusForbidden)
	}

	storage.CreateUser(username, password)
	u, err := storage.GetParamsFromDB(username, password)
	if err != nil {
		_, err := fmt.Fprintf(w, "%v %v %v", u.Username, u.Password, "Some errors")
		if err != nil {
			return
		}
	}
	log.Println("Somebody want to register 2")
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/home.html")
	if err != nil {
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}

	log.Println("Somebody want to go home")
}

func main() {
	err := storage.CreateDB()

	if err != nil {
		log.Fatal(err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				LoginFront(w, r)
			}
		case http.MethodPost:
			{
				LoginBack(w, r)
			}
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				Home(w, r)
			}
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				RegisterFront(w, r)
			}
		case http.MethodPost:
			{
				RegisterBack(w, r)
			}
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521},
	}

	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Listening on :443...")
	log.Fatal(srv.ListenAndServeTLS("./certs/server.crt", "./certs/keypair.pem"))
}
