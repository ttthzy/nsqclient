package controller

import (
	"log"
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	User := "hello"
	t, err := template.ParseFiles("template/html/home.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, User)
}
