package mypage

import (
	"app/contexthandler"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
)

var mypageTpl *template.Template

type tplProcess struct {
	data map[string]interface{}
	user model.User
	err  error
}

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

func processCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}, tName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		mypageTpl.ExecuteTemplate(w, tName, data)
		return
	}
	data[c.Name] = string(b64Str)
	mypageTpl.ExecuteTemplate(w, tName, data)
}

func ShowMypage(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	t.data["userName"] = t.user.UserName
	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.data, "mypage.html")
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "mypage.html")
		return
	}
	mypageTpl.ExecuteTemplate(w, "mypage.html", t.data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	titleNames := getRouteTitles(t.user.ID)
	t.data["userName"] = t.user.UserName
	t.data["titles"] = titleNames

	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.data, "show_routes.html")
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "show_routes.html")
		return
	}

	mypageTpl.ExecuteTemplate(w, "show_routes.html", t.data)
}

func ConfirmDelete(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.data["title"] = req.FormValue("title")
	mypageTpl.ExecuteTemplate(w, "confirm_delete.html", t.data)
}

func ShowQuestionForm(w http.ResponseWriter, req *http.Request) {
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
	mypageTpl.ExecuteTemplate(w, "question_form.html", t.data)
}

func ShowAllSimulRoutes(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	titleNames := getSimulRouteTitles(t.user.ID)
	t.data["userName"] = t.user.UserName
	t.data["titles"] = titleNames

	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.data, "show_simul_routes.html")
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "show_simul_routes.html")
		return
	}

	mypageTpl.ExecuteTemplate(w, "show_simul_routes.html", t.data)
}
