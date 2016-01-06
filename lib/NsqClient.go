package lib

import "nsqclient/models"

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

type Handle struct {
	Msgchan    chan *nsq.Message
	ChanSwitch bool
	Nci        models.Messages
}

//给客户端推送的消息
var RevMsg [2]string

var HH *Handle

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
		case <-time.After(time.Second):
			if h.ChanSwitch {
				close(h.Msgchan)
				HH = nil
				fmt.Println("关闭了")
				return
			}
		}
	}
}

func (h *Handle) Stop() {
	h.ChanSwitch = true
	close(h.Msgchan)
	HH = nil
}

func (h *Handle) SetHH() {
	HH = h
}

func Connect_Nsq(constr string, nci models.Messages) {

	// if HH.Nci.ClientID != "" {
	// 	return
	// }

	config := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(nci.Topic, nci.Channel, config)
	if err != nil {
		//panic(err)
		return
	}
	h := new(Handle)
	consumer.AddHandler(nsq.HandlerFunc(h.HandleMsg))
	h.Msgchan = make(chan *nsq.Message, 1024)
	err = consumer.ConnectToNSQD(constr)
	if err != nil {
		return
		//这里需要加一个循环计次的逻辑处理，？次以后不再尝试连接。
		//fmt.Println("连接服务器失败，尝试再次连接中...")
		//Connect_Nsq(constr, nci)
	}

	nci.ClientID = config.ClientID
	nci.UserID = GetGuid()
	h.Nci = nci
	h.SetHH()
	h.Process()
}

///接收channel消息并处理
func (h *Handle) ReceiveMessage() {
	//llogger.Info(message)
	fmt.Println("Message：" + h.Nci.Message)
	h.Nci.SendDate = time.Now()
	//ret := models.AddMessages(h.nci)
	RevMsg[0] = h.Nci.MessageID
	RevMsg[1] = h.Nci.Message
	//RevMsg = h.Nci.MessageID + "|" + h.Nci.Message
}
