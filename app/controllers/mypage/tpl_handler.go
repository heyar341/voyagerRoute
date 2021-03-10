package mypage

import (
	"app/cookiehandler"
	"app/customerr"
	"app/tplutil"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
)

var mypageTpl *template.Template

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

func processCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		mypageTpl.ExecuteTemplate(w, "show_routes.html", data)
		return
	}
	data[c.Name] = string(b64Str)
	mypageTpl.ExecuteTemplate(w, "show_routes.html", data)
}

func ShowMypage(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	t.Data["userName"] = t.User.UserName
	c, err := req.Cookie("msg")
	if err != nil {
		mypageTpl.ExecuteTemplate(w, "mypage.html", t.Data)
		return
	}
	t.Data["msg"] = c.Value
	mypageTpl.ExecuteTemplate(w, "mypage.html", t.Data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	titleNames := routeTitles(t.User.ID)
	t.Data["userName"] = t.User.UserName
	t.Data["titles"] = titleNames

	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.Data)
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.Data)
		return
	}

	mypageTpl.ExecuteTemplate(w, "show_routes.html", t.Data)
}

func ConfirmDelete(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.Data["title"] = req.FormValue("title")
	mypageTpl.ExecuteTemplate(w, "confirm_delete.html", t.Data)
}

func ShowQuestionForm(w http.ResponseWriter, req *http.Request) {
	t := tplutil.GetTplData(req)
	if t.Err != nil {
		e := t.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.Data["userName"] = t.User.UserName
	t.Data["email"] = t.User.Email
	mypageTpl.ExecuteTemplate(w, "question_form.html", t.Data)
}
