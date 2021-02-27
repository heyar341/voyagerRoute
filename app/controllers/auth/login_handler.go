package auth

import (
	"app/dbhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
)

func Login(w http.ResponseWriter, req *http.Request) {
	//エラーメッセージを定義
	msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}
	//Validation完了後のパスワードを取得
	password, ok := req.Context().Value("password").(string)
	if !ok {
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}

	//取得するドキュメントの条件
	userDoc := bson.D{{"email", email}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	if err != nil {
		msg = "メールアドレスまたはパスワードが正しくありません。"
		//入力されたメールアドレスを保持する
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		return
	}
	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		log.Printf("Error while bson marshaling user document: %v", err)
		return
	}
	var user model.UserData
	//marshalした値をUnmarshalして、userに代入
	err = bson.Unmarshal(bsonByte, &user)
	if err != nil {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		log.Printf("Error while bson unmarshaling :%v", err)
		return
	}

	storedPass := user.Password
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(password))
	//一致しない場合
	if err != nil {
		msg = "メールアドレスまたはパスワードが正しくありません。"
		//入力されたメールアドレスを保持する
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		return
	}
	//一致した場合
	err = genNewSession(user.ID, w)
	if err != nil {
		msg = " ログインに失敗しました。もう一度操作をしなおしてください。"
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
