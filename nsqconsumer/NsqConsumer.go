package main

import (
	"fmt"
	"nsqclient/lib"
	"nsqclient/models"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	NSQCONSTR string = "nsq-ttthzygi35.tenxcloud.net:40255"
)

type NsqMsgHandler struct {
	msg       models.Messages
	waitGroup sync.WaitGroup
}

///NSQ消息处理
func (handler *NsqMsgHandler) HandleMessage(nsqMsg *nsq.Message) error {
	handler.waitGroup.Add(1)
	defer handler.waitGroup.Done()

	handler.msg.Message = lib.DecodeStr(string(nsqMsg.Body))
	handler.msg.MessageID = string(nsqMsg.ID[:])
	handler.msg.SendDate = time.Now()
	fmt.Println("Got msg:", handler.msg.MessageID, handler.msg.Message)
	
    return saveMsg(handler.msg)
}

func saveMsg(msg models.Messages) error{
    fileHelper:=new(lib.FileHelper)
    
    
    return fileHelper.WriteFile("data.txt",msg.Message)
}

func main() {
	nci := models.Messages{
		Topic:   "test",
		Channel: "Jessehua",
		UserID:  "00001",
	}

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(nci.Topic, nci.Channel, config)
	if err != nil {
		fmt.Println("NSQ消费者创建失败！")
		return
	}

	handler := new(NsqMsgHandler)
	handler.waitGroup.Add(1)

	consumer.AddHandler(nsq.HandlerFunc(handler.HandleMessage))
	conErr := consumer.ConnectToNSQD(NSQCONSTR)
	if conErr != nil {
		fmt.Println("NSQ服务器连接失败")
		return
	}

	handler.waitGroup.Wait()
}
