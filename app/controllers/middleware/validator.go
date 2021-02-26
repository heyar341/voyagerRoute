package middleware

import (
	"context"
	"net/http"
	"net/url"
)

func LoginValidator(next http.HandlerFunc) http.HandlerFunc {
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
		ctx := req.Context()
		ctx = context.WithValue(ctx, "email", email)
		ctx = context.WithValue(ctx, "password", password)

		next.ServeHTTP(w, req.WithContext(ctx))
	}
}
