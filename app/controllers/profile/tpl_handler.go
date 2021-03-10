package profile

import (
	"app/cookiehandler"
	"app/customerr"
	"app/tplutil"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
)

var profileTpl *template.Template

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
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.Data["userName"] = t.User.UserName
	t.Data["email"] = t.User.Email
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.Data, "profile.html")
		return
	}
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.Data, "profile.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "profile.html", t.Data)
}

func EditUserNameForm(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.Data["userName"] = t.User.UserName
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.Data, "username_edit.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "username_edit.html", t.Data)
}
func EditEmailForm(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.Data["email"] = t.User.Email
	newEmail := req.URL.Query().Get("newEmail")
	t.Data["newEmail"] = newEmail

	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.Data, "email_edit.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "email_edit.html", t.Data)
}

func EditPasswordForm(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.Data, "password_edit.html")
		return
	}
	profileTpl.ExecuteTemplate(w, "password_edit.html", t.Data)
}
