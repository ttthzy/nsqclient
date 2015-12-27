package lib

import (
	"fmt"
	"github.com/lvanneo/llog/llogger"
	"github.com/nsqio/go-nsq"
	"time"
)

var NsqLock bool

type Handle struct {
	msgchan chan *nsq.Message
	stop    bool
	nci     NsqConnInfo
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
			ReceiveMessage(string(m.Body), h.nci)
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

func Connect_Nsq(nci NsqConnInfo) {
	NsqLock = false
	config := nsq.NewConfig()
	config.ClientID = nci.UserID
	consumer, err := nsq.NewConsumer(nci.Topic, nci.Channel, config)
	if err != nil {
		//panic(err)
		fmt.Println(err.Error())
		return
	}
	h := new(Handle)
	consumer.AddHandler(nsq.HandlerFunc(h.HandleMsg))
	h.msgchan = make(chan *nsq.Message, 1024)
	err = consumer.ConnectToNSQD("nsq-ttthzygi35.tenxcloud.net:40255")
	if err != nil {
		//这里需要加一个循环计次的逻辑处理，？次以后不再尝试连接。
		fmt.Println("连接服务器失败，尝试再次连接中...")
		Connect_Nsq(nci)
	}
	h.nci = nci
	go h.Stopcmd()
	h.Process()

}

func (h *Handle) Stopcmd() {
	for {
		if NsqLock {
			h.stop = true
			close(h.msgchan)
			fmt.Println("与服务器的连接已经关闭了")
		}
	}
}

func Cmdstp() {
	NsqLock = true
	fmt.Println("正在尝试关闭与服务器的连接通道")
}

///接收channel消息并处理
func ReceiveMessage(msg string, nci NsqConnInfo) {
	message := "jsondata={Topic:" + nci.Topic + ",Channel:" + nci.Channel + ",UserID:" + nci.UserID + ",Message:" + msg + "}"
	llogger.Info(message)
	fmt.Println(msg)
}
