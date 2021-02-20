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
	//Cookieが設定されてない場合
	if err != nil {
		c = &http.Cookie{
			Name:  "sessionId",
			Value: "",
		}
	}

	sesId, err := ParseToken(c.Value)
	if err != nil {
		msg := "ログインしていません"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Println(err)
		return
	}

	if sesId != "" {
		//DBから読み込み
		sesDoc := bson.D{{"session_id", sesId}}
		err = dbhandler.Delete("googroutes", "sessions", sesDoc)
		if err != nil {
			msg := "エラ〜が発生しました。"
			http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
			return
		}
	}

	c.MaxAge = -1
	http.SetCookie(w, c)
	msg := "ログアウトしました"
	http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)

}
