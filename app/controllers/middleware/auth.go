package middleware

import (
	"app/controllers/auth"
	"app/dbhandler"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		//templateに渡すデータ
		data := map[string]interface{}{"isLoggedIn": false}
		//Cookieからセッション情報取得
		c, err := req.Cookie("sessionId")
		//Cookieが設定されてない場合
		if err != nil {
			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		//tokenからsessionID取得
		sessionID, _ := auth.ParseToken(c.Value)
		userDoc := bson.D{{"session_id", sessionID}}
		//DBから読み込み
		_, err = dbhandler.Find("googroutes", "sessions", userDoc, nil)
		if err != nil {
			ctx = context.WithValue(ctx, "data", data)
			next.ServeHTTP(w, req.WithContext(ctx))
			return
		}

		data["isLoggedIn"] = true
		ctx = context.WithValue(ctx, "data", data)
		next.ServeHTTP(w, req.WithContext(ctx))
	}
}
