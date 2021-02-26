package auth

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	if err != nil {
		msg := "ログインしていません"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Println(err)
		return
	}

	sessionID, err := ParseToken(c.Value)
	if err != nil {
		msg := "ログインしていません"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Println(err)
		return
	}

	//DBから読み込み
	sesDoc := bson.D{{"session_id", sessionID}}
	err = dbhandler.Delete("googroutes", "sessions", sesDoc)
	if err != nil {
		msg := "ログアウト中エラーが発生しました。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	//Cookieを削除
	c.MaxAge = -1
	http.SetCookie(w, c)

	success := "ログアウトしました"
	http.Redirect(w, req, "/?success="+success, http.StatusSeeOther)
}
