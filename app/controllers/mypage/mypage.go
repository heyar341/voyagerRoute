package mypage

import (
	"app/controllers"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/model"
	"encoding/base64"
	"fmt"
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

var fieldNameMap map[string]string = map[string]string{
	"simul_search": "simul_route_titles",
	"multi_search": "multi_route_titles",
}

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

func (m *mypageController) getAndSetSearchType(req *http.Request) string {
	searchType := req.URL.Query().Get("search_type")
	switch searchType {
	case "multi_search":
		m.data["searchType"] = searchType
	case "simul_search":
		m.data["searchType"] = searchType
	default:
		m.Err = customerr.BaseErr{
			Op:  "get search type from query parameter",
			Msg: "不正なURLです。",
			Err: fmt.Errorf("error while getting search type from query parameter"),
		}
		return ""
	}
	return searchType
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

func (m *mypageController) existsCookie(w http.ResponseWriter, req *http.Request, tName string) bool {
	//successメッセージがある場合
	c, _ := req.Cookie("success")
	if c != nil {
		showMsgWithCookie(w, c, m.data, tName)
		return true
	}
	//エラーメッセージがある場合
	c, _ = req.Cookie("msg")
	if c != nil {
		showMsgWithCookie(w, c, m.data, tName)
		return true
	}
	return false
}

func Mypage(w http.ResponseWriter, req *http.Request) {
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	exsitsCookie := m.existsCookie(w, req, "mypage.html")
	if exsitsCookie {
		return
	}
	m.data["userName"] = m.user.UserName
	mypageTpl.ExecuteTemplate(w, "mypage.html", m.data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	m.getAndSetSearchType(req)
	searchType := m.getAndSetSearchType(req)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	var titleNames []string
	titleNames = getRouteTitles(m.user.ID, fieldNameMap[searchType])
	if searchType == "multi_search" {
		m.data["cardColor"] = "#1A73E8"
	} else if searchType == "simul_search" {
		m.data["cardColor"] = "#0dcaf0"
	}

	m.data["userName"] = m.user.UserName
	m.data["titles"] = titleNames
	exsitsCookie := m.existsCookie(w, req, "show_routes.html")
	if exsitsCookie {
		return
	}

	mypageTpl.ExecuteTemplate(w, "show_routes.html", m.data)
}

func ConfirmDelete(w http.ResponseWriter, req *http.Request) {
	var m mypageController
	m.data = m.GetLoginStateFromCtx(req)
	m.GetUserFromCtx(req, &m.user)
	m.getAndSetSearchType(req)
	searchType := m.getAndSetSearchType(req)
	if m.Err != nil {
		e := m.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage/show_routes/"+"?search_type="+searchType)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	m.data["title"] = req.FormValue("title")
	mypageTpl.ExecuteTemplate(w, "confirm_delete.html", m.data)
}

func QuestionForm(w http.ResponseWriter, req *http.Request) {
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
