package mypage

import (
	"app/model"
	"html/template"
	"log"
	"net/http"
)

var mypageTpl *template.Template

//エラーメッセージ
var msg string

func init() {
	mypageTpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
}

func ShowMypage(w http.ResponseWriter, req *http.Request) {
	data,ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v",ok)
		return
	}
	user,ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v",ok)
		return
	}
	data["userName"] = user.UserName
	mypageTpl.ExecuteTemplate(w, "mypage.html", data)
}

func ShowAllRoutes(w http.ResponseWriter, req *http.Request) {
	data,ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v",ok)
		return
	}
	user,ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v",ok)
		return
	}
	titleNames := RouteTitles(user.ID)
	data["userName"] = user.UserName
	data["titles"] = titleNames
	//ルートの確認画面に飛べないエラーが発生した場合用
	data["msg"] = req.URL.Query().Get("msg")
	mypageTpl.ExecuteTemplate(w, "show_routes.html", data)
}
