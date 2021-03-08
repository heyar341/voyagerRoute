package mypage

import (
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var mypageTpl *template.Template

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

type mypageData struct {
	data map[string]interface{}
	user model.UserData
	err  error
}

func getDataFromCtx(req *http.Request) *mypageData {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		return &mypageData{
			err: customerr.BaseErr{
				Op:  "Getting data from context",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while getting data from context"),
			},
		}
	}
	return &mypageData{
		data: data,
	}
}

func (m *mypageData) getUserFromCtx(req *http.Request) {
	if m.err != nil {
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		m.err = customerr.BaseErr{
			Op:  "Getting user from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from context"),
		}
		return
	}
	m.user = user
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
	m := getDataFromCtx(req)
	m.getUserFromCtx(req)
	if m.err != nil {
		e := m.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	m.data["userName"] = m.user.UserName
	c, err := req.Cookie("msg")
	if err != nil {
		mypageTpl.ExecuteTemplate(w, "mypage.html", m.data)
		return
	}
	m.data["msg"] = c.Value
	mypageTpl.ExecuteTemplate(w, "mypage.html", m.data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	m := getDataFromCtx(req)
	m.getUserFromCtx(req)
	if m.err != nil {
		e := m.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	titleNames := routeTitles(m.user.ID)
	m.data["userName"] = m.user.UserName
	m.data["titles"] = titleNames

	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, m.data)
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, m.data)
		return
	}

	mypageTpl.ExecuteTemplate(w, "show_routes.html", m.data)
}

func ConfirmDelete(w http.ResponseWriter, req *http.Request) {
	m := getDataFromCtx(req)
	if m.err != nil {
		e := m.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	m.data["title"] = req.FormValue("title")
	mypageTpl.ExecuteTemplate(w, "confirm_delete.html", m.data)
}

func ShowQuestionForm(w http.ResponseWriter, req *http.Request) {
	m := getDataFromCtx(req)
	m.getUserFromCtx(req)
	if m.err != nil {
		e := m.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	m.data["userName"] = m.user.UserName
	m.data["email"] = m.user.Email
	mypageTpl.ExecuteTemplate(w, "question_form.html", m.data)
}
