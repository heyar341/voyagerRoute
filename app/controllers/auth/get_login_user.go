package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"app/dbhandler"

)

type SessionData struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	SessionId string `json:"session_id" bson:"session_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
}

func GetLoginUserID(sessionId string) (primitive.ObjectID, error) {
	sesDoc := bson.D{{"session_id", sessionId}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "sessions", sesDoc)
	if err != nil {
		return primitive.NilObjectID, err
	}
	//DBから取得した値をmarshal
	bsonByte,err := bson.Marshal(resp)
	if err != nil {
		return primitive.NilObjectID, err
	}
	var sesData SessionData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &sesData)

	return sesData.UserId, nil
}