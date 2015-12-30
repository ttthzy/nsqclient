package controller

import (
	//"encoding/json"
	"io"
	"net/http"
	"nsqclient/lib"
	"nsqclient/models"

	"github.com/pquerna/ffjson/ffjson"
)

func GetNsqHandler(w http.ResponseWriter, req *http.Request) {
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
	req.ParseForm()
	param1, found1 := req.Form["topic"]
	result := NewBaseJsonBean()

	if !found1 {
		result.Code = -99
		result.Message = "请勿非法访问"
	}

	topic := param1[0]

	if topic != "" {

		dblist := models.GetMessagesForField("topic", topic)
		msgjson, _ := ffjson.Marshal(dblist)
		result.Code = 100
		result.Message = string(msgjson)
	} else {
		result.Code = 101
		result.Message = "topic is null"
	}

	//向客户端返回JSON数据,用到了ffjson包，据说比自带的json效率高3倍
	//bytes, _ := json.Marshal(result)
	bytes, _ := ffjson.Marshal(result)
	io.WriteString(w, string(bytes))
}
