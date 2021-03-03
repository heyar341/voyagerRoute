package middleware

import (
	"app/controllers/auth"
	"app/dbhandler"
	"app/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		//templateに渡すデータ
		data := map[string]interface{}{"isLoggedIn": false} //初期値はログインしていないとしてfalse

		//sessionデータを取得
		userID, err := getUserIDFromSession(req)
		if err != nil {
			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		//ログインユーザーのデータ取得
		user, err := getLoginUser(userID)
		if err != nil {
			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		//エラーがなければ、ログインしている
		data["isLoggedIn"] = true
		//contextに各フィールドの値を追加
		ctx = context.WithValue(ctx, "data", data)
		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, req.WithContext(ctx))
	}
}

func getUserIDFromSession(req *http.Request) (primitive.ObjectID, error) {
	//Cookieからセッション情報取得
	c, err := req.Cookie("session_id")
	//Cookieが設定されてない場合
	if err != nil {
		return primitive.NilObjectID, err
	}

	//tokenからsessionID取得
	sessionID, _ := auth.ParseToken(c.Value)
	sessionDoc := bson.D{{"session_id", sessionID}}
	//DBから読み込み
	sBSON, err := dbhandler.Find("googroutes", "sessions", sessionDoc, nil)
	if err != nil {
		log.Printf("Error while finding session data: %v", err)
		return primitive.NilObjectID, err
	}

	userID := sBSON["user_id"].(primitive.ObjectID)

	return userID, nil
}

func getLoginUser(userID primitive.ObjectID) (model.UserData, error) {
	userDoc := bson.D{{"_id", userID}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	if err != nil {
		log.Printf("Error while finding user data: %v", err)
		return model.UserData{}, err
	}
	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		log.Printf("Error while bson marshaling user data: %v", err)
		return model.UserData{}, err
	}

	var user model.UserData
	//marshalした値をUnmarshalして、userに代入
	err = bson.Unmarshal(bsonByte, &user)
	if err != nil {
		log.Printf("Error while bson unmarshaling user data: %v", err)
		return model.UserData{}, err
	}
	return user, nil
}
