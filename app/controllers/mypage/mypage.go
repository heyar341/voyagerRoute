package mypage

import (
	"app/controllers"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/model"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//user documentのmulti_route_titlesフィールドの値を入れるstruct
type titleMap struct {
	titleName string
	timeStamp time.Time
}

type titleSlice []titleMap

func (t titleSlice) Len() int      { return len(t) }
func (t titleSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

//TimestampのAfterメソッドで、ソート時に最新のタイトルが先頭に来るようにする
func (t titleSlice) Less(i, j int) bool { return t[i].timeStamp.After(t[j].timeStamp) }

func getRouteTitles(userID primitive.ObjectID, fieldName string) []string {
	b, err := model.FindUser("_id", userID)
	if err != nil {
		return []string{}
	}
	titlesM, ok := b[fieldName].(primitive.M) //bson M型 (map[string]interface{})
	if !ok {
		return []string{}
	}

	var titles = make(titleSlice, len(titlesM))
	i := 0
	for title, tStamp := range titlesM {
		t, ok := tStamp.(primitive.DateTime)
		if !ok {
			log.Println("Assertion error at checking timestamp type")
			return []string{}
		}
		timeStamp := t.Time() //time.Time型に変換
		titles[i] = titleMap{titleName: title, timeStamp: timeStamp}
		i++
	}
	//保存日時順にソート
	sort.Sort(titles)
	//タイトル名を入れるsliceを作成
	titleNames := make([]string, 0, len(titles))
	for _, tMap := range titles {
		titleNames = append(titleNames, tMap.titleName)
	}

	return titleNames
}

var mypageTpl *template.Template

type mypageController struct {
	controllers.Controller
	data map[string]interface{}
	user model.User
}

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

func showMsgWithCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}, tName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		mypageTpl.ExecuteTemplate(w, tName, data)
		return
	}
	data[c.Name] = string(b64Str)
	mypageTpl.ExecuteTemplate(w, tName, data)
}

func ShowMypage(w http.ResponseWriter, req *http.Request) {
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	m.data["userName"] = m.user.UserName
	//successメッセージがある場合
	c, _ := req.Cookie("success")
	if c != nil {
		showMsgWithCookie(w, c, m.data, "mypage.html")
		return
	}
	//エラーメッセージがある場合
	c, _ = req.Cookie("msg")
	if c != nil {
		showMsgWithCookie(w, c, m.data, "mypage.html")
		return
	}
	mypageTpl.ExecuteTemplate(w, "mypage.html", m.data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	titleType := req.URL.Query().Get("search_type")
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	var titleNames []string
	if titleType == "multi_search" {
		titleNames = getRouteTitles(m.user.ID, "multi_route_titles")
		m.data["searchType"] = "multi_search"
		m.data["cardColor"] = "#1A73E8"
	} else if titleType == "simul_search" {
		titleNames = getRouteTitles(m.user.ID, "simul_route_titles")
		m.data["searchType"] = "simul_search"
		m.data["cardColor"] = "#0dcaf0"
	} else {
		http.Error(w, "不正なURLです。", http.StatusBadRequest)
		return
	}

	m.data["userName"] = m.user.UserName
	m.data["titles"] = titleNames

	//successメッセージがある場合
	c, err := req.Cookie("success")
	if err == nil {
		showMsgWithCookie(w, c, m.data, "show_routes.html")
		return
	}
	//エラーメッセージがある場合
	c, err = req.Cookie("msg")
	if err == nil {
		showMsgWithCookie(w, c, m.data, "show_routes.html")
		return
	}
	mypageTpl.ExecuteTemplate(w, "show_routes.html", m.data)
}

func ConfirmDelete(w http.ResponseWriter, req *http.Request) {
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	m.data["title"] = req.FormValue("title")
	mypageTpl.ExecuteTemplate(w, "confirm_delete.html", m.data)
}

func ShowQuestionForm(w http.ResponseWriter, req *http.Request) {
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	m.data["userName"] = m.user.UserName
	m.data["email"] = m.user.Email
	mypageTpl.ExecuteTemplate(w, "question_form.html", m.data)
}
