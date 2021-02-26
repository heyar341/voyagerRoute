package auth

import (
	"app/dbhandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type UserData struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id"`
	UserName         string               `json:"username" bson:"username"`
	Email            string               `json:"email" bson:"email"`
	Password         []byte               `json:"password" bson:"password"`
	MultiRouteTitles map[string]time.Time `json:"multi_route_titles" bson:"multi_route_titles"`
}

var login_tpl *template.Template

func init() {
	login_tpl = template.Must(template.Must(template.ParseGlob("templates/auth/*")).ParseGlob("templates/includes/*.html"))
}

func Login(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := IsLoggedIn(req)
	if req.Method == "GET" {
		data := map[string]interface{}{"isLoggedIn": isLoggedIn, "msg": ""}
		register_tpl.ExecuteTemplate(w, "login.html", data)
		return
	}

	//メールアドレスをリクエストから取得
	email := req.FormValue("email")
	if email == "" {
		msg := "メールアドレスを入力してください。"
		data := map[string]interface{}{"isLoggedIn": isLoggedIn, "msg": msg}
		register_tpl.ExecuteTemplate(w, "login.html", data)
		return
	}
	//パスワードをリクエストから取得
	password := req.FormValue("password")
	if password == "" {
		msg := "パスワードを入力してください。"
		data := map[string]interface{}{"isLoggedIn": isLoggedIn, "msg": msg}
		register_tpl.ExecuteTemplate(w, "login.html", data)
		return
	}
	//取得するドキュメントの条件
	emailDoc := bson.D{{"email", email}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "users", emailDoc, nil)
	if err != nil {
		msg := "メールアドレスまたはパスワードが正しくありません。"
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
	}
	//DBから取得した値をmarshal
	bsonByte, _ := bson.Marshal(resp)
	var user UserData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	storedPass := user.Password
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(password))
	//一致しない場合
	if err != nil {
		msg := "メールアドレスまたはパスワードが正しくありません。"
		data := map[string]interface{}{"isLoggedIn": isLoggedIn, "msg": msg}
		register_tpl.ExecuteTemplate(w, "login.html", data)
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
		data := map[string]interface{}{"isLoggedIn": isLoggedIn, "msg": msg}
		register_tpl.ExecuteTemplate(w, "login.html", data)
		log.Println(err)
		return
	}
	signedStr, err := createToken(sessionID)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		data := map[string]interface{}{"isLoggedIn": isLoggedIn, "msg": msg}
		register_tpl.ExecuteTemplate(w, "login.html", data)
		log.Println(err)
		return
	}

	//Cookieの設定
	c := &http.Cookie{
		Name:  "sessionId",
		Value: signedStr,
		Path:  "/",
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
