package main

import (
	"nsqclient/controller"
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
	models.InitUserConsumer()
	StartHttpServer()

	//test-123
	return

	running := true
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入指令：\n")
	for running {
		data, _, _ := reader.ReadLine()
		///执行自定义的cmd命令
		switch command := string(data); command {
		case "http":
			StartHttpServer()
			running = false
		default:
			fmt.Printf("错误指令，请重新输入：\n")
		}
	}
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
	http.HandleFunc("/GetMsgDB/", controller.GetMsgForMongoDBHandler)

	http.HandleFunc("/ReceiveMsg/", controller.RevMsgHandler)
	http.HandleFunc("/ConMsq/", controller.ConMsqHandler)
	http.HandleFunc("/StopConsumer/", controller.StopConsumerHandler)

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
