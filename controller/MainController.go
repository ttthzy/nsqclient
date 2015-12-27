package controller

import (
	"encoding/json"
	"io"
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
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/html/home.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}
func GetNsqHandler(w http.ResponseWriter, req *http.Request) {
	//获取客户端通过GET/POST方式传递的参数
	req.ParseForm()
	param1, found1 := req.Form["contype"]
	result := NewBaseJsonBean()

	if !found1 {
		result.Code = -99
		result.Message = "请勿非法访问"
	}

	contype := param1[0]

	if contype == "start" {
		result.Code = 100
		result.Message = "与服务器连接成功"
		//lib.Connect_Nsq()
	} else {
		result.Code = 101
		result.Message = "与服务器断开"
		//lib.Cmdstp()
	}

	//向客户端返回JSON数据
	bytes, _ := json.Marshal(result)
	io.WriteString(w, string(bytes))

}
