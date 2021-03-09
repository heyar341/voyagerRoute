package auth

import (
	"app/model"
	"log"
	"net/http"
	"net/url"
)

func Logout(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	msg := url.QueryEscape("ログインしていません")
	//Cookieからセッション情報取得
	c, err := req.Cookie("session_id")
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Println(err)
		return
	}

	sessionID, err := ParseToken(c.Value)
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Println(err)
		return
	}

	err = model.DeleteSession(sessionID)
	if err != nil {
		log.Println(err)
		return
	}

	//Cookieを削除
	c.MaxAge = -1
	http.SetCookie(w, c)

	success := url.QueryEscape("ログアウトしました")
	http.Redirect(w, req, "/?success="+success, http.StatusSeeOther)
}
