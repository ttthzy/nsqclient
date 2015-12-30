package main

import "nsqclient/lib"
import "nsqclient/models"
import "nsqclient/controller"
import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	//StartNsq()
	//StartHttpServer()
	//return

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

///nsq客户端服务
func StartNsq() {
	nci := models.Messages{
		Topic:   "test",
		Channel: "Eason",
		UserID:  "00001",
	}
	constr := "nsq-ttthzygi35.tenxcloud.net:40255"
	lib.Connect_Nsq(constr, nci)
}

///http服务器
func StartHttpServer() {
	fmt.Printf("HttpServer Run...\n")

	//静态目录
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))

	//页面路由
	http.HandleFunc("/index.html", controller.IndexHandler)
	http.HandleFunc("/home.html", controller.HomeHandler)
	http.HandleFunc("/", controller.NotFoundHandler)

	///API接口路由
	http.HandleFunc("/SendMsg/", controller.GetNsqHandler)
	http.HandleFunc("/ReceiveMsg/", controller.GetMsgHandler)

	///启动监听服务
	http.ListenAndServe(":8080", nil)

}
