package controller

import (
	"log"
	"net/http"
	"text/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
