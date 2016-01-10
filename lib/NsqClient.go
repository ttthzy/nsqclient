package lib

import (
	"io/ioutil"
	"net/http"
	"nsqclient/models"
	"strings"
)

import (
	"fmt"
	"time"

	"github.com/lvanneo/llog/llogger"
	"github.com/nsqio/go-nsq"
)

type Handle struct {
	Msgchan    chan *nsq.Message
	ChanSwitch bool
	Nci        models.Messages
}

///推送返回的消息体
var RevMsg map[string]interface{} = make(map[string]interface{})
var UD models.UserDynamic

func (h *Handle) HandleMsg(m *nsq.Message) error {
	if !h.ChanSwitch {
		h.Msgchan <- m
	}
	return nil
}

func (h *Handle) Process() {

	h.ChanSwitch = false

	for {
		select {
		case m := <-h.Msgchan:
			h.Nci.Message = string(m.Body)
			h.Nci.MessageID = string(m.ID[:])
			h.ReceiveMessage()
		case <-time.After(2 * time.Second):
			if h.ChanSwitch {
				close(h.Msgchan)
				fmt.Println("超时")
				return
			}
			// default:
			// 	return
		}
	}
}

func (h *Handle) Stop() {
	h.ChanSwitch = true
	close(h.Msgchan)
}

func Connect_Nsq(constr string, ud models.UserDynamic) string {
	config := nsq.NewConfig()

	if GetUserNsqState(ud.UserID, config.ClientID) {
		return "connected"
	}
	consumer, err := nsq.NewConsumer(ud.Topic, ud.Channel, config)
	if err != nil {
		return err.Error()
	}

	h := new(Handle)
	consumer.AddHandler(nsq.HandlerFunc(h.HandleMsg))

	h.Msgchan = make(chan *nsq.Message, 1024)
	err = consumer.ConnectToNSQD(constr)
	if err != nil {
		return err.Error()
	}

	///记录用户状态
	ud.CreateDate = time.Now()
	ud.HostID = config.ClientID
	ud.IsOnline = true

	consumerid := models.AddUserDynamic(ud)

	h.Nci = models.Messages{
		ConsumerID: consumerid,
	}

	UD = ud
	go StopConsumer(consumer, ud.UserID, config.ClientID)
	h.Process()

	return config.ClientID

}

func StopConsumer(consumer *nsq.Consumer, UserID, ClientID string) {
	limiter := time.Tick(10 * time.Second) //设置for循环间隔时间 10秒
	for {
		<-limiter
		if !GetUserNsqState(UserID, ClientID) {
			llogger.Info(ClientID + ":stop")
			consumer.Stop()
			break
		} else {
			llogger.Info(ClientID + ":runing")
		}
	}
}

///接收channel消息并处理
func (h *Handle) ReceiveMessage() {
	h.Nci.Message = DecodeStr(h.Nci.Message)
	h.Nci.SendDate = time.Now()
	if h.Nci.MessageID == "" {
		return
	}
	if h.Nci.MessageID == "undefined" {
		return
	}
	if h.Nci.Message == "" {
		return
	}
	if h.Nci.Message == "undefined" {
		return
	}

	wtime := h.Nci.SendDate.Format("2006-01-02 15:04:05")
	//llogger.Info(message)
	fmt.Println("Message："+h.Nci.Message, "(", wtime, ")")

	models.AddMessages(h.Nci)

	RevMsg["ConsumerID"] = h.Nci.ConsumerID
	RevMsg["MssageID"] = h.Nci.MessageID
	RevMsg["Mssage"] = h.Nci.Message
	RevMsg["DateTime"] = wtime

}

///查询并返回用户nsq登录状态
func GetUserNsqState(userid, hostid string) bool {
	udb := models.GetUserDynamicByhostid(userid, hostid)
	return udb.IsOnline
}

///向nsq服务器推送一条消息
func HttpDo_NSQ(faction, furl, fdata string) string {
	client := &http.Client{}
	req, err := http.NewRequest(faction, furl, strings.NewReader(EncodeStr(fdata)))
	if err != nil {
		return "接口错误"
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "发送失败"
	}
	return string(body)
}
