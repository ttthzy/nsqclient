package lib

import "nsqclient/models"

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

type Handle struct {
	msgchan chan *nsq.Message
	stop    bool
	nci     models.Messages
}

func (h *Handle) HandleMsg(m *nsq.Message) error {
	if !h.stop {
		h.msgchan <- m
	}
	return nil
}

func (h *Handle) Process() {
	h.stop = false
	for {
		select {
		case m := <-h.msgchan:
			//fmt.Println(string(m.Body))
			h.nci.Message = string(m.Body)
			ReceiveMessage(h.nci)
		case <-time.After(time.Second):
			if h.stop {
				close(h.msgchan)
				fmt.Println("关闭了")
				return
			}
		}
	}
}

func (h *Handle) Stop() {
	h.stop = true

}

func Connect_Nsq(constr string, nci models.Messages) {
	config := nsq.NewConfig()
	config.ClientID = nci.UserID
	consumer, err := nsq.NewConsumer(nci.Topic, nci.Channel, config)
	if err != nil {
		panic(err)
		return
	}
	h := new(Handle)
	consumer.AddHandler(nsq.HandlerFunc(h.HandleMsg))
	h.msgchan = make(chan *nsq.Message, 1024)
	err = consumer.ConnectToNSQD(constr)
	if err != nil {
		//这里需要加一个循环计次的逻辑处理，？次以后不再尝试连接。
		fmt.Println("连接服务器失败，尝试再次连接中...")
		Connect_Nsq(constr, nci)
	}
	h.nci = nci
	h.Process()

}

///接收channel消息并处理
func ReceiveMessage(nci models.Messages) {
	//message := "jsondata={Topic:" + nci.Topic + ",Channel:" + nci.Channel + ",UserID:" + nci.UserID + ",Message:" + nci.Message + "}"
	//llogger.Info(message)
	nci.SendDate = time.Now()
	ret := models.AddMessages(nci)

	fmt.Println("ID：" + ret + " Message：" + nci.Message)
}