package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"app/dbhandler"

)

type SessionData struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	SessionId string `json:"sessionid" bson:"sessionid"`
	UserId primitive.ObjectID `json:"userid" bson:"userid"`
}

func GetLoginUserID(sessionId string) (primitive.ObjectID, error) {
	//DBから読み込み
	client, ctx, err := dbhandler.Connect()
	if err != nil {
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database("googroutes")
	sessionsCollection := database.Collection("sessions")
	//DBからのレスポンスを挿入する変数
	var sesData SessionData
	err = sessionsCollection.FindOne(ctx, bson.D{{"sessionid", sessionId}}).Decode(&sesData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Fatal("ドキュメントが見つかりません")
		}
		log.Fatal(err)
		return primitive.NilObjectID, err
	}
	return sesData.UserId, nil
}