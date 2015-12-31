package controller

import (
	//"encoding/json"
	"io"
	"net/http"
	"nsqclient/lib"
	"nsqclient/models"

	"github.com/pquerna/ffjson/ffjson"
	"strconv"
)


func StartNsqHandler(w http.ResponseWriter, req *http.Request) {
    
    
    
    
    
    
}





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
	var retmsg string
	req.ParseForm()
	param1, found1 := req.Form["topic"]
    param2, _ := req.Form["limit"]


	if !found1 {
		retmsg = "请勿非法访问"
	}

	topic := param1[0]
    limit,err := strconv.Atoi(param2[0])
    if err != nil{
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
