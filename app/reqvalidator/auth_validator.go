package reqvalidator

import (
	"app/cookiehandler"
	"app/customerr"
	"app/mailhandler"
	"context"
	"fmt"
	"log"
	"net/http"
)

type authValidator struct {
	username string
	email    string
	password string
	err      error
}

func (a *authValidator) checkHTTPMethod(req *http.Request) {
	if req.Method != "POST" {
		a.err = customerr.BaseErr{
			Op:  "checking HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method access"),
		}
	}
}

func (a *authValidator) getUserName(req *http.Request) {
	if a.err != nil {
		return
	}
	userName := req.FormValue("username")
	if userName == "" {
		a.err = customerr.BaseErr{
			Op:  "get username from request form",
			Msg: "ユーザー名を入力してください。",
			Err: fmt.Errorf("username was empty"),
		}
	}
	a.username = userName
}

func (a *authValidator) getEmail(req *http.Request) {
	if a.err != nil {
		return
	}
	email := req.FormValue("email")
	if email == "" {
		a.err = customerr.BaseErr{
			Op:  "get email from request form",
			Msg: "メールアドレスを入力してください。",
			Err: fmt.Errorf("email was empty"),
		}
	}
	a.email = email
}

func (a *authValidator) getPassword(req *http.Request) {
	if a.err != nil {
		return
	}
	password := req.FormValue("password")
	if password == "" {
		a.err = customerr.BaseErr{
			Op:  "get password from request form",
			Msg: "パスワードを入力してください。",
			Err: fmt.Errorf("password was empty"),
		}
	} else if len(password) < 8 {
		a.err = customerr.BaseErr{
			Op:  "get password from request form",
			Msg: "パスワードは８文字以上入力してください。",
			Err: fmt.Errorf("password was empty"),
		}
	}
	a.password = password
}

//正規表現によるメールアドレスの形式チェック、およびアドレスドメインの有効性チェックを行う
func (a *authValidator) checkEmail(email string) {
	if a.err != nil {
		return
	}
	if !mailhandler.IsEmailValid(email) {
		a.err = customerr.BaseErr{
			Op:  "check validity email syntax and domain",
			Msg: "メールアドレスに不備があります。",
			Err: fmt.Errorf("invalid email"),
		}
	}
}

func RegisterValidator(Register http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var a authValidator
		a.checkHTTPMethod(req)
		a.getUserName(req)
		a.getEmail(req)
		a.checkEmail(a.email)
		a.getPassword(req)

		if a.err != nil {
			e := a.err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/register_form")
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "username", a.username)
		ctx = context.WithValue(ctx, "email", a.email)
		ctx = context.WithValue(ctx, "password", a.password)

		Register.ServeHTTP(w, req.WithContext(ctx))
	}
}

func LoginValidator(Login http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var a authValidator
		a.checkHTTPMethod(req)
		a.getEmail(req)
		a.getPassword(req)

		if a.err != nil {
			e := a.err.(customerr.BaseErr)
			cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/login_form")
			log.Printf("operation: %s, error: %v", e.Op, e.Err)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "email", a.email)
		ctx = context.WithValue(ctx, "password", a.password)

		Login.ServeHTTP(w, req.WithContext(ctx))
	}
}
