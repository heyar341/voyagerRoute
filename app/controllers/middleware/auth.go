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
		session, err := getSession(req)
		if err != nil {
			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		//ログインユーザーのデータ取得
		user, err := getLoginUser(session.UserID)
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

func getSession(req *http.Request) (model.SessionData, error) {
	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		return model.SessionData{}, err
	}

	//tokenからsessionID取得
	sessionID, _ := auth.ParseToken(c.Value)
	sessionDoc := bson.D{{"session_id", sessionID}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "sessions", sessionDoc, nil)
	if err != nil {
		log.Printf("Error while finding session data: %v", err)
		return model.SessionData{}, err
	}

	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		log.Printf("Error while bson marshaling session data: %v", err)
		return model.SessionData{}, err
	}

	var session model.SessionData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &session)
	return session, nil
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
		log.Printf("Error while bson marshaling session data: %v", err)
		return model.UserData{}, err
	}

	var user model.UserData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	return user, nil
}
