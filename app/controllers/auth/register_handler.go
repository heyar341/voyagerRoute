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
	//DBに保存
	client, ctx, err := dbhandler.Connect()
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database("googroutes")
	usersCollection := database.Collection("users")
	insRes, err := usersCollection.InsertOne(ctx,bson.D{
		{"username",uName},
		{"password",securedPass},
	})
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	userDocId := insRes.InsertedID
	//固有のセッションIDを作成
	sesId := uuid.New().String()
	//DBに保存
	sessionsCollection := database.Collection("sessions")
	_, err = sessionsCollection.InsertOne(ctx,bson.D{
		{"session_id",sesId},
		{"user_id",userDocId},
	})
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
