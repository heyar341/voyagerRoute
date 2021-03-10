package auth

import (
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

type logoutProcess struct {
	err error
}

//getCookieFromRequest gets Cookie contains sessionID from request
func (l *logoutProcess) getCookieFromRequest(req *http.Request) *http.Cookie {
	if l.err != nil {
		return nil
	}
	c, err := req.Cookie("session_id")
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "get cookie contains sessionID form request",
			Msg: "ログイン情報が取得できません。",
			Err: fmt.Errorf("err while getting cookie from request"),
		}
		return nil
	}
	return c
}

//parseCookieToken parse token of sessionID in Cookie
func (l *logoutProcess) parseCookieToken(c *http.Cookie) string {
	if l.err != nil {
		return ""
	}
	sessionID, err := ParseToken(c.Value)
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "get sessionID from Cookie",
			Msg: "ログイン情報が取得できません。",
			Err: fmt.Errorf("err while getting sessinID from Cookie"),
		}
		return ""
	}
	return sessionID
}

//deleteSession deletes session document from sessions collection
func (l *logoutProcess) deleteSession(sessionID string) {
	if l.err != nil {
		return
	}
	err := model.DeleteSession(sessionID)
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "delete session document from sessions collection",
			Msg: "ログアウトできませんでした。",
			Err: fmt.Errorf("err while deleting session document from sessions collection: %w", err),
		}
		return
	}
}

func Logout(w http.ResponseWriter, req *http.Request) {
	var l logoutProcess
	controllers.CheckHTTPMethod(req, &l.err)
	c := l.getCookieFromRequest(req)
	sessionID := l.parseCookieToken(c)
	l.deleteSession(sessionID)

	if l.err != nil {
		e := l.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	//Cookieを削除
	c.MaxAge = -1
	http.SetCookie(w, c)

	cookiehandler.MakeCookieAndRedirect(w, req, "success", "ログアウトしました", "/")
}
