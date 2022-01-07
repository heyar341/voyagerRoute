package middleware

import (
	"app/bsonconv"
	"app/controllers/auth"
	"app/model"
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

	d, err := model.FindSession(sessionID)
	if err != nil {
		log.Printf("Error while finding session data: %v", err)
		return primitive.NilObjectID, err
	}

	userID := d["user_id"].(primitive.ObjectID)

	return userID, nil
}

func getLoginUser(userID primitive.ObjectID) (model.User, error) {
	//DBから読み込み
	d, err := model.FindUser("_id", userID)
	if err != nil {
		log.Printf("Error while finding user data: %v", err)
		return model.User{}, err
	}
	var user model.User
	var e error
	bsonconv.DocToStruct(d, &user, &e, "user")
	return user, nil
}
