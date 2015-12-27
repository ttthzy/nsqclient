package main

import "nsqclient/lib/nsq"
import "nsqclient/controller"
import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	running := true
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入指令：\n")
	for running {
		data, _, _ := reader.ReadLine()
		///执行自定义的cmd命令
		switch command := string(data); command {
		case "start nsq":
			StartNsq()
			running = false
		case "start httpserver":
			StartHttpServer()
			running = false
		default:
			fmt.Printf("错误指令，请重新输入：\n")
		}

	}
}

///启动nsq客户端连接
func StartNsq() {
	nci := lib.NsqConnInfo{
		Topic:   "test",
		Channel: "eason",
		UserID:  "00001",
	}
	lib.Connect_Nsq(nci)
}

///启动http服务器
func StartHttpServer() {
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/index.html", controller.IndexHandler)
	http.HandleFunc("/home.html", controller.HomeHandler)
	http.HandleFunc("/getnsq/", controller.GetNsqHandler)
	http.HandleFunc("/", controller.NotFoundHandler)
	http.ListenAndServe(":8080", nil)

	fmt.Printf("HttpServer Run...\n")
}
