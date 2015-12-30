package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"nsqclient/lib/mongodb"
)

type Messages struct {
	Id      bson.ObjectId `bson:"_id"`
	Topic   string        `bson:"topic"`   //
	Channel string        `bson:"channel"` //
	UserID  string        `bson:"userid"`  //
	Message string        `bson:"message"` //
}

/**
 * 添加Messages对象
 */
func AddMessages(p Messages) string {
	p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := lib.WitchCollection("Messages", query)
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
	err := lib.WitchCollection("Messages", exop)
	if err != nil {
		return "true"
	}
	return "false"
}

//获取所有的person数据
func PageMessages() []Messages {
	var list []Messages
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&list)
	}
	err := lib.WitchCollection("Messages", query)
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
	lib.WitchCollection("Messages", query)
	return item
}
