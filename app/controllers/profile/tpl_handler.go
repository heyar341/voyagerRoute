package profile

import (
	"app/model"
	"html/template"
	"log"
	"net/http"
)

var profileTpl *template.Template

//エラーメッセージ
var msg string

func init() {
	profileTpl = template.Must(template.Must(template.ParseGlob("templates/profile/*.html")).ParseGlob("templates/includes/*.html"))
}

func ShowProfile(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	data["userName"] = user.UserName
	data["email"] = user.Email

	profileTpl.ExecuteTemplate(w, "profile.html", data)
}

func EditUserNameForm(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	data["userName"] = user.UserName

	profileTpl.ExecuteTemplate(w, "username_edit.html", data)
}
func EditEmailForm(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	data["email"] = user.Email

	profileTpl.ExecuteTemplate(w, "email_edit.html", data)
}
func EditPasswordForm(w http.ResponseWriter, req *http.Request) {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}

	profileTpl.ExecuteTemplate(w, "password_edit.html", data)
}
