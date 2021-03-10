package profile

import (
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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
	controllers.GetUserFromCtx(req, &uP.user, &uP.err)
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
