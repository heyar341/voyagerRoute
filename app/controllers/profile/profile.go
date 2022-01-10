package profile

import (
	"app/controllers"
	"app/internal/bsonconv"
	"app/internal/contexthandler"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/internal/mailhandler"
	"app/model"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type updateUserName struct {
	user        model.User
	newUserName string
	err         error
}

const REDIRECT_URI_TO_UPDATE_USERNAME_FORM = "/profile/username_edit_form"

//getUserNameFromForm gets username from request form.
func (uU *updateUserName) getUserNameFromForm(req *http.Request) {
	if uU.err != nil {
		return
	}
	newUserName := req.FormValue("username")

	if newUserName == "" {
		uU.err = customerr.BaseErr{
			Op:  "get username from request form",
			Msg: "ユーザー名は１文字以上入力してください。",
			Err: fmt.Errorf("request's username was empty"),
		}
		return
	}
	uU.newUserName = newUserName
}

//updateUserName updates user document's username field.
func (uU *updateUserName) updateUserName() {
	if uU.err != nil {
		return
	}
	err := model.UpdateUser(uU.user.ID, "username", uU.newUserName)
	if err != nil {
		uU.err = customerr.BaseErr{
			Op:  "Saving editing email to DB",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting editing email to editing_email collection %w", err),
		}
		return
	}
}

func UpdateUserName(w http.ResponseWriter, req *http.Request) {
	var uU updateUserName
	controllers.CheckHTTPMethod(req, &uU.err)
	contexthandler.GetUserFromCtx(req, &uU.user, &uU.err)
	uU.getUserNameFromForm(req)
	uU.updateUserName()

	if uU.err != nil {
		e := uU.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_USERNAME_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	cookiehandler.MakeCookieAndRedirect(w, req, "success", "ユーザー名を変更しました。", "/profile")
}

const REDIRECT_URI_TO_UPDATE_EMAIL_FORM = "/profile/email_edit_form"

type updateEmailProcess struct {
	user     model.User
	newEmail string
	err      error
}

//getEmailFromForm gets email from request's form
func (u *updateEmailProcess) getEmailFromForm(req *http.Request) {
	if u.err != nil {
		return
	}
	newEmail := req.FormValue("email")
	if !mailhandler.IsEmailValid(newEmail) {
		u.err = customerr.BaseErr{
			Op:  "check email address's validity",
			Msg: "メールアドレスに不備があります。",
			Err: fmt.Errorf("request email was invalid %v", newEmail),
		}
		return
	}
	u.newEmail = newEmail
}

//saveEditingEmailToDB saves editing email to DB
func (u *updateEmailProcess) saveEditingEmailToDB(token string) {
	if u.err != nil {
		return
	}
	err := model.SaveEditingEmail(u.newEmail, token)
	if err != nil {
		u.err = customerr.BaseErr{
			Op:  "Saving editing email to DB",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting editing email to editing_email collection %w", err),
		}
		return
	}
}

func UpdateEmail(w http.ResponseWriter, req *http.Request) {
	var u updateEmailProcess
	controllers.CheckHTTPMethod(req, &u.err)
	contexthandler.GetUserFromCtx(req, &u.user, &u.err)
	u.getEmailFromForm(req)
	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()
	u.saveEditingEmailToDB(token)

	if u.err != nil {
		e := u.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_EMAIL_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	//認証依頼画面表示
	http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

	//メールでトークン付きのURLを送る
	err := mailhandler.SendConfirmEmail(token, u.newEmail, "confirm_email")
	return
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}

}

type confirmUpdateEmail struct {
	user         model.User
	editingEmail model.EditingEmail
	token        string
	err          error
}

//getTokenFromURL gets token from URL parameter
func (c *confirmUpdateEmail) getTokenFromURL(req *http.Request) {
	if c.err != nil {
		return
	}
	token := req.URL.Query()["token"][0]
	if token == "" {
		c.err = customerr.BaseErr{
			Op:  "get token form URL",
			Msg: "トークン情報が不正です。",
			Err: fmt.Errorf("token was empty"),
		}
	}
	c.token = token
}

//getEditingEmailDocFromDB fetch editing email from DB
func (c *confirmUpdateEmail) getEditingEmailDocFromDB() bson.M {
	if c.err != nil {
		return nil
	}
	d, err := model.GetEditingEmailDoc(c.token)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			c.err = customerr.BaseErr{
				Op:  "get editing email document form DB",
				Msg: "認証トークンが一致しません。",
				Err: fmt.Errorf("token was invalid: %w", err),
			}
		default:
			c.err = customerr.BaseErr{
				Op:  "get editing email document form DB",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while finding editing email document from editing_email collection %w", err),
			}
		}
		return nil
	}
	return d
}

//checkTokenExpire checks if token expires or not
func (c *confirmUpdateEmail) checkTokenExpire() {
	if c.err != nil {
		return
	}
	t := time.Unix(c.editingEmail.ExpiresAt, 0)
	if !t.After(time.Now()) {
		c.err = customerr.BaseErr{
			Op:  "check if token expires or not",
			Msg: "トークンの有効期限が切れています。もう一度メールアドレス変更のお手続きをしてください。",
			Err: fmt.Errorf("token expired"),
		}
		return
	}
}

//updateUserEmail updates email field in user document
func (c *confirmUpdateEmail) updateUserEmail() {
	if c.err != nil {
		return
	}
	err := model.UpdateUser(c.user.ID, "email", c.editingEmail.Email)
	if err != nil {
		c.err = customerr.BaseErr{
			Op:  "update use document's email field",
			Msg: "メールアドレスの更新に失敗しました。",
			Err: fmt.Errorf("err while updating email field in user document: %w", err),
		}
		return
	}
}

func ConfirmUpdateEmail(w http.ResponseWriter, req *http.Request) {
	var c confirmUpdateEmail
	controllers.CheckHTTPMethod(req, &c.err)
	contexthandler.GetUserFromCtx(req, &c.user, &c.err)
	c.getTokenFromURL(req)
	d := c.getEditingEmailDocFromDB()
	bsonconv.DocToStruct(d, &c.editingEmail, &c.err, "editing email")
	c.checkTokenExpire()
	c.updateUserEmail()

	if c.err != nil {
		e := c.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_EMAIL_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	cookiehandler.MakeCookieAndRedirect(w, req, "success¬", "メールアドレスの変更が完了しました。", "/profile")
}

const REDIRECT_URI_TO_UPDATE_PASSWORD_FORM = "/profile/password_edit_form"

type updatePassword struct {
	user            model.User
	securedPassword []byte
	err             error
}

//getPasswordFromForm gets password from request form.
func (uP *updatePassword) getPasswordFromForm(req *http.Request, fieldName string) string {
	if uP.err != nil {
		return ""
	}
	p := req.FormValue(fieldName)
	if len(p) < 8 {
		uP.err = customerr.BaseErr{
			Op:  "get password from request form",
			Msg: "パスワードは８文字以上入力してください。",
			Err: fmt.Errorf("request's pasword length was invalid"),
		}
		return ""
	}
	return p
}

//comparePasswords compares hashed password in DB and password user inputted
func (uP *updatePassword) comparePasswords(p string) {
	if uP.err != nil {
		return
	}
	err := bcrypt.CompareHashAndPassword(uP.user.Password, []byte(p))
	if err != nil {
		uP.err = customerr.BaseErr{
			Op:  "compare passwords",
			Msg: "パスワードが間違っています。",
			Err: fmt.Errorf("password was not right: %w", err),
		}
		return
	}
}

//hashPassword hashes password user inputted
func (uP *updatePassword) hashPassword(p string) {
	if uP.err != nil {
		return
	}
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(p), 5)
	if err != nil {
		uP.err = customerr.BaseErr{
			Op:  "hash password",
			Msg: "パスワードの更新に失敗しました。",
			Err: fmt.Errorf("error while hashing password: %w", err),
		}
		return
	}
	uP.securedPassword = securedPassword
}

//updatePassword updates user document's password field
func (uP *updatePassword) updatePassword() {
	if uP.err != nil {
		return
	}
	err := model.UpdateUser(uP.user.ID, "password", uP.securedPassword)
	if err != nil {
		uP.err = customerr.BaseErr{
			Op:  "update password",
			Msg: "パスワードの更新に失敗しました。",
			Err: fmt.Errorf("error while updating password: %w", err),
		}
		return
	}
}

func UpdatePassword(w http.ResponseWriter, req *http.Request) {
	var uP updatePassword
	controllers.CheckHTTPMethod(req, &uP.err)
	contexthandler.GetUserFromCtx(req, &uP.user, &uP.err)
	currPassword := uP.getPasswordFromForm(req, "current-password")
	newPassword := uP.getPasswordFromForm(req, "password")
	uP.comparePasswords(currPassword)
	uP.hashPassword(newPassword)
	uP.updatePassword()

	if uP.err != nil {
		e := uP.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_PASSWORD_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	cookiehandler.DeleteCookie(w, "session_id", "/")

	success := "パスワードの変更に成功しました。ログインしてください。"
	cookiehandler.MakeCookieAndRedirect(w, req, "success", success, "/login_form")
}

var profileTpl *template.Template

type tplProcess struct {
	data map[string]interface{}
	user model.User
	err  error
}

func init() {
	profileTpl = template.Must(template.Must(template.ParseGlob("templates/profile/*.html")).ParseGlob("templates/includes/*.html"))
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
	c, err := req.Cookie("success")
	if err == nil {
		processCookie(w, c, t.data, "profile.html")
		return
	}
	c, err = req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "profile.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "profile.html", t.data)
}

func EditUserNameForm(w http.ResponseWriter, req *http.Request) {
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
	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "username_edit.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "username_edit.html", t.data)
}
func EditEmailForm(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	t.data["email"] = t.user.Email
	newEmail := req.URL.Query().Get("newEmail")
	t.data["newEmail"] = newEmail

	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "email_edit.html")
		return
	}

	profileTpl.ExecuteTemplate(w, "email_edit.html", t.data)
}

func EditPasswordForm(w http.ResponseWriter, req *http.Request) {
	var t tplProcess
	t.data = contexthandler.GetLoginStateFromCtx(req)
	contexthandler.GetUserFromCtx(req, &t.user, &t.err)
	if t.err != nil {
		e := t.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/mypage")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	c, err := req.Cookie("msg")
	if err == nil {
		processCookie(w, c, t.data, "password_edit.html")
		return
	}
	profileTpl.ExecuteTemplate(w, "password_edit.html", t.data)
}
