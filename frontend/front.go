package frontend

import (
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

	log.Println("Somebody want to login")
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
