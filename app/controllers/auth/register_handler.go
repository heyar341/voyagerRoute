package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"net/url"
	"golang.org/x/crypto/bcrypt"
	"app/dbhandler"
	"github.com/google/uuid"
)

func Register(w http.ResponseWriter, req *http.Request){
	if req.Method != "POST"{
		msg := url.QueryEscape("HTTPメソッドが不正です。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//ユーザー名をリクエストから取得
	uName := req.FormValue("username")
	if uName == ""{
		msg := url.QueryEscape("ユーザー名を入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをリクエストから取得
	password := req.FormValue("password")
	if password == ""{
		msg := url.QueryEscape("パスワードを入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをハッシュ化
	securedPass,err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	//userをDBに保存
	//保存するドキュメント
	userDoc := bson.D{
		{"username",uName},
		{"password",securedPass},
	}
	//DBに保存
	insertRes, err := dbhandler.Insert("googroutes", "users", userDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	//insertResから、userのドキュメントIDを取得
	userDocId := insertRes.InsertedID
	//固有のセッションIDを作成
	sesId := uuid.New().String()
	//sessionをDBに保存
	sesDoc := bson.D{
		{"session_id",sesId},
		{"user_id",userDocId},
	}
	_, err = dbhandler.Insert("googroutes", "sessions", sesDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	signedStr,err := createToken(sesId)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	//Cookieの設定
	c := &http.Cookie{
		Name: "sessionId",
		Value: signedStr,
	}
	http.SetCookie(w,c)
	http.Redirect(w,req,"/",http.StatusSeeOther)
}
