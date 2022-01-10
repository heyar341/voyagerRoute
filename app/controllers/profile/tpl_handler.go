package profile

import (
	"app/internal/contexthandler"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/model"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
)

var profileTpl *template.Template

type tplProcess struct {
	data map[string]interface{}
	user model.User
	err  error
}

func init() {
	profileTpl = template.Must(template.Must(template.ParseGlob("templates/profile/*.html")).ParseGlob("templates/includes/*.html"))
}

func processCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}, tplName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		profileTpl.ExecuteTemplate(w, tplName, data)
		return
	}
	data[c.Name] = string(b64Str)
	profileTpl.ExecuteTemplate(w, tplName, data)
}

func ShowProfile(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.data["userName"] = t.user.UserName
	t.data["email"] = t.user.Email
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.data, "profile.html")
		return
	}
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "profile.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "profile.html", t.data)
}

func EditUserNameForm(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.data["userName"] = t.user.UserName
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "username_edit.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "username_edit.html", t.data)
}
func EditEmailForm(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.data["email"] = t.user.Email
	newEmail := req.URL.Query().Get("newEmail")
	t.data["newEmail"] = newEmail

	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "email_edit.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "email_edit.html", t.data)
}

func EditPasswordForm(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "password_edit.html")
		return
	}
	profileTpl.ExecuteTemplate(w, "password_edit.html", t.data)
}
