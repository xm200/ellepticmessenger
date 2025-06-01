package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func generateKeyPair(w http.ResponseWriter, r *http.Request) error {
	serverCurve := ecdh.X25519()
	PrivateKey, err := serverCurve.GenerateKey(rand.Reader)
	if err != nil {
		_, err := fmt.Fprintf(w, "Error occured, try later %v", err)
		if err != nil {
			return err
		}
		return err
	}

	PrivKey := ""
	for _, b := range PrivateKey.Bytes() {
		PrivKey = PrivKey + fmt.Sprintf("%x", b)
	}

	PubKey := ""
	for _, b := range PrivateKey.PublicKey().Bytes() {
		PubKey = PubKey + fmt.Sprintf("%x", b)
	}

	_, err = fmt.Fprintf(w, "{\"PrivKey\":\"%v\",\"Pubkey\":\"%v\"}", PrivKey, PubKey)
	if err != nil {
		return err
	}

	log.Println("Generated keypair")
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Will be right back")
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Will be right back")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				Login(w, r)
			}
		case http.MethodPost:
			{
				Login(w, r)
			}
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				Register(w, r)
			}
		case http.MethodPost:
			{
				Register(w, r)
			}
		default:
			http.Error(w, "Method now allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/keypair", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			generateKeyPair(w, r)
		} else {
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
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Listening on :443...")
	log.Fatal(srv.ListenAndServeTLS("./certs/server.crt", "./certs/keypair.pem"))
}
