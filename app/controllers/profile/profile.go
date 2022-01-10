package profile

import (
	"app/controllers"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/internal/mailhandler"
	"app/model"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type profileController struct {
	controllers.Controller
	user        model.User
	updateValue string
}

var profileTpl *template.Template

func init() {
	profileTpl = template.Must(template.Must(template.ParseGlob("templates/profile/*.html")).ParseGlob("templates/includes/*.html"))
}

const redirectURIToUpdateUsernameForm = "/profile/username_edit"

//getUserNameFromForm gets username from request form.
func (p *profileController) getUserNameFromForm(req *http.Request) {
	if p.Err != nil {
		return
	}
	newUserName := req.FormValue("username")

	if newUserName == "" {
		p.Err = customerr.BaseErr{
			Op:  "get username from request form",
			Msg: "ユーザー名は１文字以上入力してください。",
			Err: fmt.Errorf("request's username was empty"),
		}
		return
	}
	p.updateValue = newUserName
}

//updateUserName updates user document's username field.
func (p *profileController) updateUserName() {
	if p.Err != nil {
		return
	}
	err := model.UpdateUser(p.user.ID, "username", p.updateValue)
	if err != nil {
		p.Err = customerr.BaseErr{
			Op:  "Saving editing email to DB",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting editing email to editing_email collection %w", err),
		}
		return
	}
}

func EditUserName(w http.ResponseWriter, req *http.Request) {
	var p profileController
	switch req.Method {
	case "GET":
		data := p.GetLoginStateFromCtx(req)
		p.GetUserFromCtx(req, &p.user)
		if p.Err != nil {
			e := p.Err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}
		data["userName"] = p.user.UserName
		c, _ := req.Cookie("msg")
		if c != nil {
			processCookie(w, c, data, "username_edit.html")
			return
		}

		profileTpl.ExecuteTemplate(w, "username_edit.html", data)

	case "POST":
		p.GetUserFromCtx(req, &p.user)
		p.getUserNameFromForm(req)
		p.updateUserName()
		if p.Err != nil {
			e := p.Err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, redirectURIToUpdateUsernameForm)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		cookiehandler.MakeCookieAndRedirect(w, req, "success", "ユーザー名を変更しました。", "/profile")

	default:
		http.Error(w, "不正なHTTPメソッドです。", http.StatusMethodNotAllowed)
	}
}

const redirectURIToUpdateEmailForm = "/profile/email_edit"

//getEmailFromForm gets email from request's form
func (p *profileController) getEmailFromForm(req *http.Request) {
	if p.Err != nil {
		return
	}
	newEmail := req.FormValue("email")
	if !mailhandler.IsEmailValid(newEmail) {
		p.Err = customerr.BaseErr{
			Op:  "check email address's validity",
			Msg: "メールアドレスに不備があります。",
			Err: fmt.Errorf("request email was invalid %v", newEmail),
		}
		return
	}
	p.updateValue = newEmail
}

//saveEditingEmailToDB saves editing email to DB
func (p *profileController) saveEditingEmailToDB(token string) {
	if p.Err != nil {
		return
	}
	err := model.SaveEditingEmail(p.updateValue, token)
	if err != nil {
		p.Err = customerr.BaseErr{
			Op:  "Saving editing email to DB",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting editing email to editing_email collection %w", err),
		}
		return
	}
}

func EditEmail(w http.ResponseWriter, req *http.Request) {
	var p profileController
	switch req.Method {
	case "GET":
		data := p.GetLoginStateFromCtx(req)
		p.GetUserFromCtx(req, &p.user)
		if p.Err != nil {
			e := p.Err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}
		data["email"] = p.user.Email
		newEmail := req.URL.Query().Get("newEmail")
		data["newEmail"] = newEmail

		c, _ := req.Cookie("msg")
		if c != nil {
			processCookie(w, c, data, "email_edit.html")
			return
		}

		profileTpl.ExecuteTemplate(w, "email_edit.html", data)

	case "POST":
		p.GetUserFromCtx(req, &p.user)
		p.getEmailFromForm(req)
		//メールアドレス認証用のトークンを作成
		token := uuid.New().String()
		p.saveEditingEmailToDB(token)

		if p.Err != nil {
			e := p.Err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, redirectURIToUpdateEmailForm)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}
		//認証依頼画面表示
		http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

		//メールでトークン付きのURLを送る
		err := mailhandler.SendConfirmEmail(token, p.updateValue, "confirm_email")
		return
		if err != nil {
			log.Printf("Error while sending email at registering: %v", err)
		}

	default:
		http.Error(w, "不正なHTTPメソッドです。", http.StatusMethodNotAllowed)
	}
}

const redirectURIToUpdatePasswordForm = "/profile/password_edit"

//getPasswordFromForm gets password from request form.
func (p *profileController) getPasswordFromForm(req *http.Request, fieldName string) string {
	if p.Err != nil {
		return ""
	}
	pwd := req.FormValue(fieldName)
	if len(pwd) < 8 {
		p.Err = customerr.BaseErr{
			Op:  "get password from request form",
			Msg: "パスワードは８文字以上入力してください。",
			Err: fmt.Errorf("request's pasword length was invalid"),
		}
		return ""
	}
	return pwd
}

//comparePasswords compares hashed password in DB and password user inputted
func (p *profileController) comparePasswords(pwd string) {
	if p.Err != nil {
		return
	}
	err := bcrypt.CompareHashAndPassword(p.user.Password, []byte(pwd))
	if err != nil {
		p.Err = customerr.BaseErr{
			Op:  "compare passwords",
			Msg: "パスワードが間違っています。",
			Err: fmt.Errorf("password was not right: %w", err),
		}
		return
	}
}

//hashPassword hashes password user inputted
func (p *profileController) hashPassword(pwd string) []byte {
	if p.Err != nil {
		return []byte{}
	}
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), 5)
	if err != nil {
		p.Err = customerr.BaseErr{
			Op:  "hash password",
			Msg: "パスワードの更新に失敗しました。",
			Err: fmt.Errorf("error while hashing password: %w", err),
		}
		return []byte{}
	}
	return securedPassword
}

//updatePassword updates user document's password field
func (p *profileController) updatePassword(securedPassword []byte) {
	if p.Err != nil {
		return
	}
	err := model.UpdateUser(p.user.ID, "password", securedPassword)
	if err != nil {
		p.Err = customerr.BaseErr{
			Op:  "update password",
			Msg: "パスワードの更新に失敗しました。",
			Err: fmt.Errorf("error while updating password: %w", err),
		}
		return
	}
}

func EditPassword(w http.ResponseWriter, req *http.Request) {
	var p profileController
	switch req.Method {
	case "GET":
		data := p.GetLoginStateFromCtx(req)
		p.GetUserFromCtx(req, &p.user)
		if p.Err != nil {
			e := p.Err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		c, _ := req.Cookie("msg")
		if c != nil {
			processCookie(w, c, data, "password_edit.html")
			return
		}
		profileTpl.ExecuteTemplate(w, "password_edit.html", data)

	case "POST":
		p.GetUserFromCtx(req, &p.user)
		currPassword := p.getPasswordFromForm(req, "current-password")
		newPassword := p.getPasswordFromForm(req, "password")
		p.comparePasswords(currPassword)
		securedPassword := p.hashPassword(newPassword)
		p.updatePassword(securedPassword)

		if p.Err != nil {
			e := p.Err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, redirectURIToUpdatePasswordForm)
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		cookiehandler.DeleteCookie(w, "session_id", "/")

		success := "パスワードの変更に成功しました。ログインしてください。"
		cookiehandler.MakeCookieAndRedirect(w, req, "success", success, "/login_form")

	default:
		http.Error(w, "不正なHTTPメソッドです。", http.StatusMethodNotAllowed)
	}
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
	var p profileController
	data := p.GetLoginStateFromCtx(req)
	p.GetUserFromCtx(req, &p.user)
	if p.Err != nil {
		e := p.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	data["userName"] = p.user.UserName
	data["email"] = p.user.Email
	c, _ := req.Cookie("success")
	if c != nil {
		processCookie(w, c, data, "profile.html")
		return
	}
	c, _ = req.Cookie("msg")
	if c != nil {
		processCookie(w, c, data, "profile.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "profile.html", data)

}
