package mypage

import (
	"app/internal/contexthandler"
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
type TitleMap struct {
	TitleName string
	TimeStamp time.Time
}

type TitileSlice []TitleMap

func (t TitileSlice) Len() int      { return len(t) }
func (t TitileSlice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

//TimestampのAfterメソッドで、ソート時に最新のタイトルが先頭に来るようにする
func (t TitileSlice) Less(i, j int) bool { return t[i].TimeStamp.After(t[j].TimeStamp) }

func getRouteTitles(userID primitive.ObjectID) []string {
	b, err := model.FindUser("_id", userID)
	if err != nil {
		return []string{}
	}
	titlesM, ok := b["multi_route_titles"].(primitive.M) //bson M型 (map[string]interface{})
	if !ok {
		return []string{}
	}

	var titles = make(map[string]time.Time)
	for title, tStamp := range titlesM {
		t, ok := tStamp.(primitive.DateTime)
		if !ok {
			log.Println("Assertion error at checking timestamp type")
			return []string{}
		}
		timeStamp := t.Time() //time.Time型に変換
		titles[title] = timeStamp
	}

	tSlice := make(TitileSlice, len(titles))
	i := 0
	for k, v := range titles {
		tSlice[i] = TitleMap{k, v}
		i++
	}
	//保存日時順にソート
	sort.Sort(tSlice)
	//タイトル名を入れるsliceを作成
	titleNames := make([]string, 0, len(titles))
	for _, tMap := range tSlice {
		titleNames = append(titleNames, tMap.TitleName)
	}

	return titleNames
}

func getSimulRouteTitles(userID primitive.ObjectID) []string {
	b, err := model.FindUser("_id", userID)
	if err != nil {
		return []string{}
	}
	titlesM, ok := b["simul_route_titles"].(primitive.M) //bson M型 (map[string]interface{})
	if !ok {
		return []string{}
	}
	var titles = make(map[string]time.Time)
	for title, tStamp := range titlesM {
		t, ok := tStamp.(primitive.DateTime)
		if !ok {
			log.Println("Assertion error at checking timestamp type")
			return []string{}
		}
		timeStamp := t.Time() //time.Time型に変換
		titles[title] = timeStamp
	}

	tSlice := make(TitileSlice, len(titles))
	i := 0
	for k, v := range titles {
		tSlice[i] = TitleMap{k, v}
		i++
	}
	//保存日時順にソート
	sort.Sort(tSlice)
	//タイトル名を入れるsliceを作成
	titleNames := make([]string, 0, len(titles))
	for _, tMap := range tSlice {
		titleNames = append(titleNames, tMap.TitleName)
	}

	return titleNames
}

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