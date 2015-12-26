package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/lvanneo/llog/llogger"
	"github.com/nsqio/go-nsq"
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
		Connect_Nsq()
	} else {
		result.Code = 101
		result.Message = "与服务器断开"
		Cmdstp()
	}

	//向客户端返回JSON数据
	bytes, _ := json.Marshal(result)
	io.WriteString(w, string(bytes))

}

func main() {
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/home.html", HomeHandler)
	http.HandleFunc("/getnsq/", GetNsqHandler)
	http.HandleFunc("/", NotFoundHandler)

	http.ListenAndServe(":8080", nil)
}

var Lock bool

type Handle struct {
	msgchan chan *nsq.Message
	stop    bool
}

func (h *Handle) HandleMsg(m *nsq.Message) error {
	if !h.stop {
		h.msgchan <- m
	}
	return nil
}

func (h *Handle) Process() {
	h.stop = false
	for {
		select {
		case m := <-h.msgchan:
			llogger.Info(string(m.Body))
			//fmt.Println(string(m.Body))
		case <-time.After(time.Second):
			if h.stop {
				close(h.msgchan)
				fmt.Println("关闭了")
				return
			}
		}
	}
}

func (h *Handle) Stop() {
	h.stop = true

}

func Connect_Nsq() {
	Lock = false
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("test", "eason", config)
	if err != nil {
		panic(err)
	}
	h := new(Handle)
	consumer.AddHandler(nsq.HandlerFunc(h.HandleMsg))
	h.msgchan = make(chan *nsq.Message, 1024)
	err = consumer.ConnectToNSQD("nsq-ttthzygi35.tenxcloud.net:40255")
	if err != nil {
		panic(err)
	}
	go h.Stopcmd()
	h.Process()

}

func (h *Handle) Stopcmd() {
	for {
		if Lock {
			h.stop = true
			close(h.msgchan)
			fmt.Println("关闭了")
		}
	}
}

func Cmdstp() {
	Lock = true
	log.Println("close")
}

////
