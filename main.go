package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net/http"
	"log"
)


func generateKeyPair(w http.ResponseWriter, r *http.Request) {
	serverCurve := ecdh.X25519()
	ClientPrivKey, err := serverCurve.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(w, "Error occured, try later")
		return
	}
	fmt.Fprintf(w, "%v, %v", ClientPrivKey, ClientPrivKey.PublicKey())
	log.Println("Generated keypair")
}


func Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("username") == "" || r.URL.Query().Get("password") == "" {
		fmt.Fprintf(w, "Provide username and password to login")
		return
	}
	log.Println("Somebody on login")
}


func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("email") == "" || r.URL.Query().Get("username") == "" || r.URL.Query().Get("password") == "" {
		fmt.Fprintf(w, "Provide email, username, password")
		return
	}
	log.Println("Somebody on reqister")
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
		MinVersion: tls.VersionTLS12,
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

