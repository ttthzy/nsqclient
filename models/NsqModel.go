package models

//用户接收到的消息体
type NsqConnInfo struct {
	Topic   string
	Channel string
	UserID  string
	Message string
}
