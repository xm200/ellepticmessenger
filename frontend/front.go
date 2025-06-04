package frontend

import (
	"emess/backend"
	"emess/storage"
	"html/template"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/register.html")
	if err != nil {
		log.Println(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Somebody want to register")
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/login.html")
	if err != nil {
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
	log.Println("Somebody want to login")
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/home.html")
	if err != nil {
		return
	}

	type page struct {
		HaveAccount bool
		Users       []storage.Nickname
	}

	account := page{}
	account.HaveAccount = backend.SessionGetter(w, r)

	if account.HaveAccount {
		account.Users = storage.OnlineGet()
	}

	err = tmpl.Execute(w, account)
	if err != nil {
		log.Println(err)
	}

	log.Println("Somebody want to go home")
}
