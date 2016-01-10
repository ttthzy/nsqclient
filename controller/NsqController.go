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

	ud := models.UserDynamic{
		Topic:   topic,
		Channel: channel,
		UserID:  userid,
	}

	ret := lib.Connect_Nsq(conadd1, ud)

	io.WriteString(w, ret)

}

///stopConsumer
func StopConsumerHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	param1, found1 := r.Form["userid"]
	param2, found2 := r.Form["hostid"]

	if !found1 || !found2 {
		io.WriteString(w, "请勿非法访问")
		return
	}

	userid := param1[0]
	hostid := param2[0]

	///更新用户状态
	// query := bson.M{"userid": userid, "hostid": hostid}
	// change := bson.M{"$set": bson.M{"isonline": false}}
	// var ret string
	// if models.UpdateUserDynamic(query, change) {
	// 	ret = "update ok"
	// } else {
	// 	ret = "update fial"
	// }

	///记录用户状态
	ud := models.UserDynamic{
		Topic:      lib.UD.Topic,
		Channel:    lib.UD.Channel,
		UserID:     userid,
		HostID:     hostid,
		IsOnline:   false,
		CreateDate: time.Now(),
	}
	ret := models.AddUserDynamic(ud)

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
