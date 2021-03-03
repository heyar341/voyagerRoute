package auth

import (
	"app/dbhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
)

func Login(w http.ResponseWriter, req *http.Request) {
	//エラーメッセージを定義
	msg := "エラーが発生しました。もう一度操作をしなおしてください。"
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		//入力されたメールアドレスを保持する
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}
	//Validation完了後のパスワードを取得
	password, ok := req.Context().Value("password").(string)
	if !ok {
		//入力されたメールアドレスを保持する
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}

	//以下のコード内のエラーメッセージ
	msg2 := "メールアドレスまたはパスワードが正しくありません。"

	//取得するドキュメントの条件
	userDoc := bson.D{{"email", email}}
	//DBから取得
	respUDoc, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	//入力されたメールアドレスを保持する
	email = url.QueryEscape(email)
	if err != nil {
		http.Redirect(w, req, "/login_form/?msg="+msg2+"&email="+email, http.StatusSeeOther)
		return
	}

	storedPassB,ok := respUDoc["password"].(primitive.Binary)
	if !ok {
		http.Redirect(w, req, "/login_form/?msg="+msg2+"&email="+email, http.StatusSeeOther)
		return
	}
	storedPass := storedPassB.Data
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(password))
	//一致しない場合
	if err != nil {
		http.Redirect(w, req, "/login_form/?msg="+msg2+"&email="+email, http.StatusSeeOther)
		return
	}
	//一致した場合
	userID := respUDoc["_id"].(primitive.ObjectID)
	err = genNewSession(userID, w)
	if err != nil {
		msg = " ログインに失敗しました。もう一度操作をしなおしてください。"
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
