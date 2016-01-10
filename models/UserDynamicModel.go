package models

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserConsumer struct {
	Id         bson.ObjectId `bson:"_id"`
	Topic      string        `bson:"topic"`      //
	Channel    string        `bson:"channel"`    //
	UserID     string        `bson:"userid"`     //
	HostID     string        `bson:"hostid"`     //
	IsOnline   bool          `bson:"isonline"`   //
	CreateDate time.Time     `bson:"createdate"` //
}

/**
 * 添加对象
 */
func AddUserConsumer(p UserConsumer) string {
	p.Id = bson.NewObjectId()
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := WitchCollection("UserConsumer", query)
	if err != nil {
		return "false"
	}
	return p.Id.Hex()
}

//更新数据
func UpdateUserConsumer(query bson.M, change bson.M) bool {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := WitchCollection("UserConsumer", exop)
	if err != nil {
		return false
	}
	return true
}

//获取所有数据
func PageUserConsumer() []UserConsumer {
	var list []UserConsumer
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&list)
	}
	err := WitchCollection("UserConsumer", query)
	if err != nil {
		return list
	}
	return list
}

//根据指定字段查询数据
func GetUserConsumerForField(fields map[string]interface{}, sort string, limit int) []UserConsumer {
	var list []UserConsumer
	qstr := make(bson.M)
	for k, v := range qstr {
		qstr[k] = v
	}
	query := func(c *mgo.Collection) error {
		return c.Find(qstr).Sort(sort).Limit(limit).All(&list)
	}
	err := WitchCollection("UserConsumer", query)
	if err != nil {
		return list
	}
	return list
}

/**
 * 获取一条记录通过objectid
 */
func GetUserConsumerById(id string) *UserConsumer {
	objid := bson.ObjectIdHex(id)
	item := new(UserConsumer)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&item)
	}
	WitchCollection("UserConsumer", query)
	return item
}

/**
 * 获取一条记录通过hostid和userid
 */
func GetUserConsumerByhostid(userid, hostid string) *UserConsumer {
	item := new(UserConsumer)
	query := func(c *mgo.Collection) error {
		qstr := bson.M{"userid": userid, "hostid": hostid}
		return c.Find(qstr).Sort("-createdate").One(&item)
	}
	WitchCollection("UserConsumer", query)
	return item
}
