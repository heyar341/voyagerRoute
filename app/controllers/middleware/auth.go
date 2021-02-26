package middleware

import (
	"app/controllers/auth"
	"app/dbhandler"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type SessionData struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	SessionId string             `json:"session_id" bson:"session_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type UserData struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	UserName         string               `json:"username" bson:"username"`
	Email            string               `json:"email" bson:"email"`
	Password         []byte               `json:"password" bson:"password"`
	MultiRouteTitles map[string]time.Time `json:"multi_route_titles" bson:"multi_route_titles"`
}

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

func getSession(req *http.Request) (SessionData, error) {
	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		return SessionData{}, err
	}

	//tokenからsessionID取得
	sessionID, _ := auth.ParseToken(c.Value)
	sessionDoc := bson.D{{"session_id", sessionID}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "sessions", sessionDoc, nil)
	if err != nil {
		log.Printf("Error while finding session data: %v", err)
		return SessionData{}, err
	}

	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		log.Printf("Error while bson marshaling session data: %v", err)
		return SessionData{}, err
	}

	var session SessionData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &session)
	return session, nil
}

func getLoginUser(userID primitive.ObjectID) (UserData, error) {
	userDoc := bson.D{{"_id", userID}}
	//DBから読み込み
	resp, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	if err != nil {
		log.Printf("Error while finding user data: %v", err)
		return UserData{}, err
	}
	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		log.Printf("Error while bson marshaling session data: %v", err)
		return UserData{}, err
	}

	var user UserData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	return user, nil
}
