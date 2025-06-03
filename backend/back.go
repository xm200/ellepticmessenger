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

type BufferedWriter struct {
	http.ResponseWriter
	buffer []byte
}

func (b *BufferedWriter) Write(data []byte) (int, error) {
	b.buffer = append(b.buffer, data...)
	return len(data), nil
}

func (b *BufferedWriter) flush() error {
	_, err := b.ResponseWriter.Write(b.buffer)
	return err
}

var jwtkey []byte

func SetSecretKey() {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Println(err)
	}
	jwtkey = key
	err := os.Setenv("SECRET_KEY", string(key))
	if err != nil {
		log.Println(err)
	}
}

func SessionSetter(username string) string {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Println("Error while signing token")
	}
	return tokenString
}

func SessionGetter(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session")

	if cookie == nil || cookie.Value == "" || err != nil {
		return false
	}

	var token *jwt.Token

	if cookie != nil && cookie.Value != "" {
		token, err = jwt.ParseWithClaims(cookie.Value, &Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SECRET_KEY")), nil
			})
	}

	if err != nil {
		log.Println("Error while parsing token")
		log.Println(err)
		return false
	}

	if token == nil {
		return false
	}

	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return true
	}

	return false
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
	if SessionGetter(w, r) {
		http.Redirect(w, r, "/", 302)
	}

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
	_, err = fmt.Fprintf(w, "Your keypair is %v %v", priv, pub)
	if err != nil {
		log.Println(err)
	}

	log.Println("Successfully logged in")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Set-Cookie", "session="+SessionSetter(username)+"; HttpOnly; Secure; SameSite=Strict")
	http.Redirect(w, r, "/home", http.StatusFound)
}

func Register(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Authentication failed", http.StatusForbidden)
	}

	if SessionGetter(w, r) {
		http.Redirect(w, r, "/", 302)
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
