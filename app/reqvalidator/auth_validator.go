package reqvalidator

import (
	"app/controllers"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/internal/mailhandler"
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

//getUserNameFromForm gets string value from request's form
func (a *authValidator) getUserNameFromForm(req *http.Request) {
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

//getEmailFromForm gets string value from request's form
func (a *authValidator) getEmailFromForm(req *http.Request) {
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

//getPasswordFromForm gets string value from request's form
func (a *authValidator) getPasswordFromForm(req *http.Request) {
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

//checkEmail checks email's format using regex and a validity of emil's domain
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
		if req.Method == "GET" {
			ctx := req.Context()
			Register.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		var a authValidator
		controllers.CheckHTTPMethod(req, &a.err)
		a.getUserNameFromForm(req)
		a.getEmailFromForm(req)
		a.checkEmail(a.email)
		a.getPasswordFromForm(req)

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
		if req.Method == "GET" {
			ctx := req.Context()
			Login.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		var a authValidator
		controllers.CheckHTTPMethod(req, &a.err)
		a.getEmailFromForm(req)
		a.getPasswordFromForm(req)

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
