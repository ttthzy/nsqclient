package controller

import (
	//"encoding/json"
	"io"
	"net/http"
	"nsqclient/lib"
	"nsqclient/models"
	"strconv"
	"time"

	"github.com/pquerna/ffjson/ffjson"
)

var (
	conadd1 = "nsq-ttthzygi35.tenxcloud.net:40255"
	conadd2 = "120.24.210.90:4150"
)

///连接nsq
func ConMsqHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	param1, found1 := r.Form["topic"]
	param2, found2 := r.Form["channel"]
	param3, found3 := r.Form["userid"]

	if !found1 || !found2 || !found3 {
		io.WriteString(w, "请勿非法访问")
		return
	}

	topic := param1[0]
	channel := param2[0]
	userid := param3[0]

	nci := models.Messages{
		Topic:   topic,
		Channel: channel,
		UserID:  userid,
	}

	go lib.Connect_Nsq(conadd1, nci)

}

///nsqInfo
func InfoMsqHandler(w http.ResponseWriter, r *http.Request) {
	h := lib.HH

	io.WriteString(w, h.Nci.UserID)

}

///ESM推送
func RevMsgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")


	m := make(map[string]string)
	m["MssageID"] = lib.RevMsg[0]
	m["Mssage"] = lib.RevMsg[1]
	m["DateTime"] = time.Now().String()

	//
	bytes, _ := ffjson.Marshal(m)
	msg := "data:" + string(bytes) + "\n\n"
	io.WriteString(w, msg)

}

func PostMsgHandler(w http.ResponseWriter, req *http.Request) {
	//获取客户端通过GET/POST方式传递的参数
	req.ParseForm()
	param1, found1 := req.Form["sendmsg"]
	result := NewBaseJsonBean()

	if !found1 {
		result.Code = -99
		result.Message = "请勿非法访问"
	}

	fdata := param1[0]

	if fdata != "" {
		faction := "POST"
		furl := "http://nsq-ttthzygi35.tenxcloud.net:20157/put?topic=test"
		result.Code = 100
		result.Message = lib.HttpDo_NSQ(faction, furl, fdata)
	} else {
		result.Code = 101
		result.Message = "消息不能为空，请重新发送！"
	}

	//向客户端返回JSON数据,用到了ffjson包，据说比自带的json效率高3倍
	//bytes, _ := json.Marshal(result)
	bytes, _ := ffjson.Marshal(result)
	io.WriteString(w, string(bytes))

}

///从数据库pull消息
func GetMsgHandler(w http.ResponseWriter, req *http.Request) {
	//获取客户端通过GET/POST方式传递的参数
	var retmsg string
	req.ParseForm()
	param1, found1 := req.Form["topic"]
	param2, _ := req.Form["limit"]

	if !found1 {
		retmsg = "请勿非法访问"
	}

	topic := param1[0]
	limit, err := strconv.Atoi(param2[0])
	if err != nil {
		retmsg = "请勿非法访问"
	}

	if topic != "" {

		dblist := models.GetMessagesForField("topic", topic, limit)
		msgjson, _ := ffjson.Marshal(dblist)
		retmsg = string(msgjson)
	} else {

		retmsg = "topic is null"
	}
	io.WriteString(w, retmsg)

}
