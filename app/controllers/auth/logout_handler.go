package auth

import (
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

type logoutErr customerr.BaseErr

//checkHTTPMethod checks HTTP method
func (l *logoutErr) checkHTTPMethod(req *http.Request) {
	if req.Method != "POST" {
		l.Op = "check HTTP method"
		l.Msg = "HTTPメソッドが不正です。"
		l.Err = fmt.Errorf("invalid HTTP method at logout")
	}
}

//getCookie gets Cookie contains sessionID from request
func (l *logoutErr) getCookie(req *http.Request) *http.Cookie {
	if l.Err != nil {
		return nil
	}
	c, err := req.Cookie("session_id")
	if err != nil {
		l.Op = "get cookie contains sessionID form request"
		l.Msg = "ログイン情報が取得できません。"
		l.Err = fmt.Errorf("err while getting cookie from request")
		return nil
	}
	return c
}

//parseCookieToken parse token of sessionID in Cookie
func (l *logoutErr) parseCookieToken(c *http.Cookie) string {
	if l.Err != nil {
		return ""
	}
	sessionID, err := ParseToken(c.Value)
	if err != nil {
		l.Op = "get sessionID from Cookie"
		l.Msg = "ログイン情報が取得できません。"
		l.Err = fmt.Errorf("err while getting sessinID from Cookie")
		return ""
	}
	return sessionID
}

//deleteSession deletes session document from sessions collection
func (l *logoutErr) deleteSession(sessionID string) {
	if l.Err != nil {
		return
	}
	err := model.DeleteSession(sessionID)
	if err != nil {
		l.Op = "delete session document from sessions collection"
		l.Msg = "ログアウトできませんでした。"
		l.Err = fmt.Errorf("err while deleting session document from sessions collection: %w", err)
		return
	}
}

func Logout(w http.ResponseWriter, req *http.Request) {
	var l logoutErr
	l.checkHTTPMethod(req)
	c := l.getCookie(req)
	sessionID := l.parseCookieToken(c)
	l.deleteSession(sessionID)

	if l.Err != nil {
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", l.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", l.Op, l.Err)
		return
	}

	//Cookieを削除
	c.MaxAge = -1
	http.SetCookie(w, c)

	cookiehandler.MakeCookieAndRedirect(w, req, "success", "ログアウトしました", "/")
}
