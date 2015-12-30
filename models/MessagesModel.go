package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Messages struct {
	Id       bson.ObjectId `bson:"_id"`
	Topic    string        `bson:"topic"`    //
	Channel  string        `bson:"channel"`  //
	UserID   string        `bson:"userid"`   //
	Message  string        `bson:"message"`  //
	SendDate time.Time     `bson:"senddate"` //
}

/**
 * 添加Messages对象
 */
func AddMessages(p Messages) string {
	p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := WitchCollection("Messages", query)
	if err != nil {
		return "false"
	}
	return p.Id.Hex()
}

//更新Messages数据
func UpdateMessages(query bson.M, change bson.M) string {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := WitchCollection("Messages", exop)
	if err != nil {
		return "true"
	}
	return "false"
}

//获取所有的Messages数据
func PageMessages() []Messages {
	var list []Messages
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&list)
	}
	err := WitchCollection("Messages", query)
	if err != nil {
		return list
	}
	return list
}

//根据指定字段查询Messages数据
func GetMessagesForField(fieldname string, fieldvalue string) []Messages {
	var list []Messages
	qstr := bson.M{fieldname: fieldvalue}
	//sort := bson.M{"senddate": "-1"}
	query := func(c *mgo.Collection) error {
		return c.Find(qstr).Sort("-senddate").Limit(1).All(&list)
	}
	err := WitchCollection("Messages", query)
	if err != nil {
		return list
	}
	return list
}

/**
 * 获取一条记录通过objectid
 */
func GetMessagesById(id string) *Messages {
	objid := bson.ObjectIdHex(id)
	item := new(Messages)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&item)
	}
	WitchCollection("Messages", query)
	return item
}
