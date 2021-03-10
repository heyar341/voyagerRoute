package model

import (
	"app/dbhandler"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type EditingEmail struct {
	Email     string `bson:"email"`
	ExpiresAt int64  `bson:"expires_at"`
	Token     string `bson:"token"`
}

func SaveEditingEmail(newEmail, token string) error {
	//保存するドキュメント
	editingDoc := bson.D{
		{"email", newEmail},
		{"expires_at", time.Now().Add(1 * time.Hour).Unix()},
		{"token", token},
	}
	//editing_email collectionに保存
	_, err := dbhandler.Insert("googroutes", "editing_email", editingDoc)
	return err
}

func GetEditingEmailDoc(token string) (bson.M, error) {
	//取得するドキュメントの条件
	tokenDoc := bson.D{{"token", token}}
	//DBから取得
	d, err := dbhandler.Find("googroutes", "editing_email", tokenDoc, nil)
	return d, err
}
