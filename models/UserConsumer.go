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

func AddUserConsumerCheckOne(p UserConsumer) string {
	user := GetUserConsumerByhostid(p.UserID, p.HostID)

	if user.Id.Hex() == "" {
		p.Id = bson.NewObjectId()
		query := func(c *mgo.Collection) error {
			return c.Insert(p)
		}
		err := WitchCollection("UserConsumer", query)
		if err != nil {
			return "failed"
		}
		return p.Id.Hex()
	} else {
		SetUserOnlineState(user.Id.Hex(), true)
		return user.Id.Hex()
	}

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

///初始化用户登录状态，全部置为false
func InitUserConsumer() bool {
	exop := func(c *mgo.Collection) error {
		query := bson.M{"isonline": true}
		change := bson.M{"$set": bson.M{"isonline": false}}
		return c.Update(query, change)
	}
	err := WitchCollection("UserConsumer", exop)
	if err != nil {
		return false
	}
	return true
}

//更新数据根据Consumerid
func UpdateUserConsumerForID(id string, update bson.M) bool {
	objid := bson.ObjectIdHex(id)
	exop := func(c *mgo.Collection) error {
		return c.UpdateId(objid, update)
	}
	err := WitchCollection("UserConsumer", exop)
	if err != nil {
		return false
	}
	return true
}

//更新数据根据Consumerid
func SetUserOnlineState(id string, state bool) bool {
	objid := bson.ObjectIdHex(id)
	change := bson.M{"$set": bson.M{"isonline": state, "createdate": time.Now()}}
	exop := func(c *mgo.Collection) error {
		return c.UpdateId(objid, change)
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
