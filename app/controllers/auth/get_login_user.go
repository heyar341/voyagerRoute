package auth

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type SessionData struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	SessionId string             `json:"session_id" bson:"session_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
}

func GetLoginUserID(req *http.Request) (primitive.ObjectID, error) {
	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		return primitive.NilObjectID, err
	}

	sessionID, _ := ParseToken(c.Value)

	userDoc := bson.D{{"session_id", sessionID}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "sessions", userDoc, nil)
	if err != nil {
		return primitive.NilObjectID, err
	}
	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		return primitive.NilObjectID, err
	}
	var user SessionData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	return user.UserID, nil
}

func GetLoginUserName(req *http.Request) (string, error) {
	userID, err := GetLoginUserID(req)
	userDoc := bson.D{{"_id", userID}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	if err != nil {
		return "", err
	}
	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		return "", err
	}
	var user UserData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	return user.UserName, nil
}
