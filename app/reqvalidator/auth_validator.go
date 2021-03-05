package reqvalidator

import (
	"app/mailhandler"
	"context"
	"net/http"
	"net/url"
)

func RegisterValidator(Register http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			msg := url.QueryEscape("HTTPメソッドが不正です。")
			http.Redirect(w, req, "/register_form/?msg="+msg, http.StatusSeeOther)
			return
		}
		//ユーザー名をリクエストから取得
		userName := req.FormValue("username")
		//メールアドレスをリクエストから取得
		email := req.FormValue("email")
		u := url.QueryEscape(userName)
		m := url.QueryEscape(email)
		if userName == "" {
			msg := url.QueryEscape("ユーザー名を入力してください。")
			http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+u+"&email="+m, http.StatusSeeOther)
			return
		}
		if email == "" {
			msg := url.QueryEscape("メールアドレスを入力してください。")
			http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+u+"&email="+m, http.StatusSeeOther)
			return
		} else if !mailhandler.IsEmailValid(email) {
			msg := url.QueryEscape("メールアドレスに不備があります。")
			http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+u+"&email="+m, http.StatusSeeOther)
			return
		}
		//パスワードをリクエストから取得
		password := req.FormValue("password")
		if password == "" {
			msg := url.QueryEscape("パスワードを入力してください。")
			http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+u+"&email="+m, http.StatusSeeOther)
			return
		} else if len(password) < 8 {
			msg := url.QueryEscape("パスワードは8文字以上で入力してください。")
			http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+u+"&email="+m, http.StatusSeeOther)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "username", userName)
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "password", password)

		Register.ServeHTTP(w, req.WithContext(ctx))
	}
}

func LoginValidator(Login http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			msg := url.QueryEscape("HTTPメソッドが不正です。")
			http.Redirect(w, req, "/login_form/?msg="+msg, http.StatusFound)
			return
		}
		//メールアドレスをリクエストから取得
		email := req.FormValue("email")
		if email == "" {
			msg := url.QueryEscape("メールアドレスを入力してください。")
			http.Redirect(w, req, "/login_form/?msg="+msg, http.StatusSeeOther)
			return
		}
		//パスワードをリクエストから取得
		password := req.FormValue("password")
		if password == "" {
			msg := url.QueryEscape("パスワードを入力してください。")
			//入力されたメールアドレスを保持する
			email = url.QueryEscape(email)
			http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "password", password)

		Login.ServeHTTP(w, req.WithContext(ctx))
	}
}
