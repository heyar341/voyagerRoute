package auth

import (
	"app/dbhandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"net/http"
	"net/url"
)

type userData struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password []byte             `json:"password" bson:"password"`
}

func Login(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		msg := url.QueryEscape("HTTPメソッドが不正です。")
		http.Redirect(w, req, "/register/?msg="+msg, http.StatusSeeOther)
		return
	}
	//メールアドレスをリクエストから取得
	email := html.EscapeString(req.FormValue("email"))
	if email == "" {
		msg := url.QueryEscape("メールアドレスを入力してください。")
		http.Redirect(w, req, "/register/?msg="+msg, http.StatusSeeOther)
		return
	}
	//パスワードをリクエストから取得
	password := html.EscapeString(req.FormValue("password"))
	if password == "" {
		msg := url.QueryEscape("パスワードを入力してください。")
		http.Redirect(w, req, "/register/?msg="+msg, http.StatusSeeOther)
		return
	}
	//取得するドキュメントの条件
	emailDoc := bson.D{{"email", email}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "users", emailDoc)
	if err != nil {
		msg := "メールアドレスまたはパスワードが正しくありません。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
	}
	//DBから取得した値をmarshal
	bsonByte, _ := bson.Marshal(resp)
	var user userData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	storedPass := user.Password
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(password))
	//一致しない場合
	if err != nil {
		msg := "メールアドレスまたはパスワードが正しくありません。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
	}
	//一致した場合
	//固有のセッションIDを作成
	sessionID := uuid.New().String()
	//sessionをDBに保存
	sessionDoc := bson.D{
		{"session_id", sessionID},
		{"user_id", user.ID},
	}
	_, err = dbhandler.Insert("googroutes", "sessions", sessionDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(err)
		return
	}
	signedStr, err := createToken(sessionID)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		log.Println(err)
		return
	}

	//Cookieの設定
	c := &http.Cookie{
		Name:  "sessionId",
		Value: signedStr,
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
