package models

import "gopkg.in/mgo.v2/bson"

type Messages struct {
	Id      bson.ObjectId `bson:"_id"`
	Topic   string        `bson:"topic"`   //
	Channel string        `bson:"channel"` //
	UserID  string        `bson:"userid"`  //
	Message string        `bson:"message"` //
}
