package controller

import (
	"log"
	"net/http"
	"text/template"
)

type BaseJsonBean struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewBaseJsonBean() *BaseJsonBean {
	return &BaseJsonBean{}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/home.html", http.StatusFound)
	}

	t, err := template.ParseFiles("template/html/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)

}
