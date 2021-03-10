package model

import (
	"app/dbhandler"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Registering struct {
	UserName  string `bson:"username"`
	Email     string `bson:"email"`
	Password  []byte `bson:"password"`
	ExpiresAt int64  `bson:"expires_at"`
	Token     string `bson:"token"`
}

func SaveRegisteringUser(userName, email, token string, securedPassword []byte) error {
	registerDoc := bson.D{
		{"username", userName},
		{"email", email},
		{"password", securedPassword},
		{"expires_at", time.Now().Add(1 * time.Hour).Unix()},
		{"token", token},
	}
	//DBに保存
	_, err := dbhandler.Insert("googroutes", "registering", registerDoc)
	return err
}

func FindUserByToken(token string) (bson.M, error) {
	//取得するドキュメントの条件
	tokenDoc := bson.M{"token": token}
	//DBから取得
	d, err := dbhandler.Find("googroutes", "registering", tokenDoc, nil)
	return d, err
}
