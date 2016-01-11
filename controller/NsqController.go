package controller

import (
	//"encoding/json"
	"io"
	"net/http"
	"nsqclient/lib"
	"nsqclient/models"
	"strconv"

	"github.com/pquerna/ffjson/ffjson"
	//"github.com/gorilla/sessions"
)

var (
	conadd1 = "nsq-ttthzygi35.tenxcloud.net:40255"
	conadd2 = "120.24.210.90:4150"
)

///声明一个session仓库
//var store = sessions.NewCookieStore([]byte("something-very-secret"))

///连接nsq
func ConMsqHandler(w http.ResponseWriter, r *http.Request) {
	values := lib.GetUrlValue(r)

	ud := models.UserConsumer{
		Topic:   values["topic"],
		Channel: values["channel"],
		UserID:  values["userid"],
	}

	ret := lib.Connect_Nsq(conadd1, values["consumerid"], ud)

	io.WriteString(w, ret)

}

///stopConsumer
func StopConsumerHandler(w http.ResponseWriter, r *http.Request) {
	values := lib.GetUrlValue(r)

	///更新用户状态

	var ret string
	if models.SetUserOnlineState(values["consumerid"], false) {
		ret = "stop ok" + values["consumerid"]
	} else {
		ret = "stop fial" + values["consumerid"]
	}

	io.WriteString(w, ret)

}

///ESM推送
func RevMsgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	//

	bytes, _ := ffjson.Marshal(lib.RevMsg)
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
func GetMsgForMongoDBHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	//获取客户端通过GET/POST方式传递的参数
	var retmsg string
	req.ParseForm()
	param1, found1 := req.Form["topic"]
	param2, _ := req.Form["sort"]
	param3, _ := req.Form["limit"]

	if !found1 {
		retmsg = "请勿非法访问"
	}

	topic := param1[0]
	sort := param2[0]
	limit, err := strconv.Atoi(param3[0])
	if err != nil {
		retmsg = "请勿非法访问"
	}

	if topic != "" {
		selectM := make(map[string]interface{})
		selectM["topic"] = topic
		selectM["isdel"] = false
		dblist := models.GetMessagesForField(selectM, sort, limit)
		msgjson, _ := ffjson.Marshal(dblist)
		retmsg = string(msgjson)
	} else {

		retmsg = "topic is null"
	}

	msg := "data:" + retmsg + "\n\n"
	io.WriteString(w, msg)

}
