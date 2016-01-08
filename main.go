package main

import (
	"nsqclient/controller"
	"nsqclient/lib"
	"nsqclient/models"
)
import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/netutil"
)

func main() {
	//StartNsq()
	StartHttpServer()

	//test
	return

	running := true
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入指令：\n")
	for running {
		data, _, _ := reader.ReadLine()
		///执行自定义的cmd命令
		switch command := string(data); command {
		case "nsq":
			StartNsq()
			running = false
		case "http":
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
	http.HandleFunc("/test.html", controller.TestHandler)
	http.HandleFunc("/", controller.NotFoundHandler)

	///API接口路由
	http.HandleFunc("/SendMsg/", controller.PostMsgHandler)
	http.HandleFunc("/GetMsg/", controller.GetMsgHandler)

	http.HandleFunc("/ReceiveMsg/", controller.RevMsgHandler)
	http.HandleFunc("/ConMsq/", controller.ConMsqHandler)
	http.HandleFunc("/InfoMsq/", controller.InfoMsqHandler)

	//http.ListenAndServe(":8080", nil)
	///启动监听服务
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Listen: %v", err)
	}
	defer l.Close()
	l = netutil.LimitListener(l, 1000000) //最大连接数

	http.Serve(l, nil)

}
