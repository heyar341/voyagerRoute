package profile

import (
	"app/model"
	"html/template"
	"log"
	"net/http"
)

var profileTpl *template.Template

func init() {
	profileTpl = template.Must(template.Must(template.ParseGlob("templates/profile/*.html")).ParseGlob("templates/includes/*.html"))
}

func ShowProfile(w http.ResponseWriter, req *http.Request) {
	msg := "エラ〜が発生しました。もう一度操作しなおしてください。"
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	data["userName"] = user.UserName
	data["email"] = user.Email
	msg = req.URL.Query().Get("msg")
	data["msg"] = msg
	success := req.URL.Query().Get("success")
	data["success"] = success
	profileTpl.ExecuteTemplate(w, "profile.html", data)
}

func EditUserNameForm(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作しなおしてください。"
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	data["userName"] = user.UserName
	msg = req.URL.Query().Get("msg")
	data["msg"] = msg

	profileTpl.ExecuteTemplate(w, "username_edit.html", data)
}
func EditEmailForm(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作しなおしてください。"
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting user from context: %v", ok)
		return
	}
	data["email"] = user.Email
	msg = req.URL.Query().Get("msg")
	data["msg"] = msg
	newEmail := req.URL.Query().Get("newEmail")
	data["newEmail"] = newEmail
	profileTpl.ExecuteTemplate(w, "email_edit.html", data)
}
func EditPasswordForm(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作しなおしてください。"
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		http.Redirect(w, req, "/mypage/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting data from context: %v", ok)
		return
	}

	msg = req.URL.Query().Get("msg")
	data["msg"] = msg
	profileTpl.ExecuteTemplate(w, "password_edit.html", data)
}
