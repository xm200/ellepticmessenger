package frontend

import (
	"emess/backend"
	"html/template"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/html/register.html")
	if err != nil {
		log.Println(err)
	}

	if !backend.SessionGetter(w, r) {
		http.Redirect(w, r, "/login", http.StatusFound)
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

	type HaveAccount struct {
		HaveAccount bool
	}

	account := HaveAccount{backend.SessionGetter(w, r)}

	err = tmpl.Execute(w, account)
	if err != nil {
		log.Println(err)
	}

	log.Println("Somebody want to go home")
}
