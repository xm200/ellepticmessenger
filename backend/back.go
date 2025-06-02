package backend

import (
	"crypto/ecdh"
	"crypto/rand"
	"emess/storage"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
	"time"
)

type Claims struct {
	Username string
	jwt.RegisteredClaims
}

var jwtkey []byte

func SetSecretKey() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Println(err)
	}
	jwtkey = key
	os.Setenv("SECRET_KEY", string(key))
}

func SessionSetter(username string) (string, string) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Println("Error while signing token")
	}
	return tokenString, expirationTime.Format(time.RFC3339)
}

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

func Login(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	log.Println(username, password)
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

	tokenString, _ := SessionSetter(username)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer:"+tokenString)
}

func Register(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/login", http.StatusFound)
	log.Println("Somebody want to register 2")
}
