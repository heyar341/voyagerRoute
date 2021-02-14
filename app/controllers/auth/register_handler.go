package auth

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"golang.org/x/crypto/bcrypt"
	"app/dbhandler"
	"github.com/google/uuid"
	"os"
	"sync"
)

func Register(w http.ResponseWriter, req *http.Request){
	if req.Method != "POST"{
		msg := url.QueryEscape("HTTPメソッドが不正です。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//ユーザー名をリクエストから取得
	userName := html.EscapeString(req.FormValue("username"))
	if userName == ""{
		msg := url.QueryEscape("ユーザー名を入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//メールアドレスをリクエストから取得
	email := html.EscapeString(req.FormValue("email"))
	if email == ""{
		msg := url.QueryEscape("メールアドレスを入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをリクエストから取得
	password := html.EscapeString(req.FormValue("password"))
	if password == ""{
		msg := url.QueryEscape("パスワードを入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	} else if len(password) < 8 {
		msg := url.QueryEscape("パスワードは8文字以上で入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをハッシュ化
	securedPassword,err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Println(err)
		return
	}
	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()

	//userを仮登録としてDBに保存
	//保存するドキュメント
	registerDoc := bson.D{
		{"username",userName},
		{"email",email},
		{"password",securedPassword},
		{"token",token},
	}
	//DBに保存
	_, err = dbhandler.Insert("googroutes", "registering", registerDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Println(err)
		return
	}

	//メール送信に少し時間がかかるので、メール送信と認証依頼画面表示を並列処理
	wg := &sync.WaitGroup{}
	wg.Add(2)

	//並列処理1:メールでトークン付きのURLを送る
	go func(){
		env_err := godotenv.Load("env/dev.env")
		if env_err != nil {
			log.Println("Can't load env file")
		}
		//envファイルからGmailのアプリパスワード取得
		gmailPassword := os.Getenv("GMAIL_APP_PASS")
		mailAuth := smtp.PlainAuth(
			"",
			"app.goog.routes@gmail.com",
			gmailPassword,
			"smtp.gmail.com",
		)

		tokenURL := "http://localhost:8080/confirm_register/?token=" + token //localhostは本番で変更
		err = smtp.SendMail(
			"smtp.gmail.com:587",
			mailAuth,
			"app.goog.routes@gmail.com",
			[]string{email},
			[]byte(fmt.Sprintf("To:%s\r\nSubject:メールアドレス認証のお願い\r\n\r\n%s", userName, tokenURL)),
		)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	//並列処理2:認証依頼画面表示
	go func() {
		http.Redirect(w,req,"/ask_confirm",http.StatusSeeOther)
		wg.Done()
	}()

	wg.Wait()
}


func ConfirmRegister(w http.ResponseWriter, req *http.Request){
	//メール認証トークンをリクエストURLから取得
	query := req.URL.Query()
	token := query["token"][0]
	if token == ""{
		http.Redirect(w,req,"/register_form",http.StatusSeeOther)
		return
	}

	//DBのregistering collectionから、user情報取得
	type registeringUser struct {
		ID primitive.ObjectID `json:"id" bson:"_id"`
		Username string `json:"username" bson:"username"`
		Email string `json:"email" bson:"email"`
		Password []byte `json:"password" bson:"password"`
		Token string `json:"token" bson:"token"`
	}

	//取得するドキュメントの条件
	tokenDoc := bson.D{{"token",token}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "registering", tokenDoc)
	if err != nil {
		msg := "認証トークンが一致しません。"
		http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)
	}
	//DBから取得した値をmarshal
	bsonByte,_ := bson.Marshal(resp)

	var user registeringUser
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	//userをDBに保存
	//保存するドキュメント
	userDoc := bson.D{
		{"username",user.Username},
		{"email",user.Email},
		{"password",user.Password},
	}
	//DBに保存
	insertRes, err := dbhandler.Insert("googroutes", "users", userDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Println(err)
		return
	}
	//insertResから、userのドキュメントIDを取得
	userDocID := insertRes.InsertedID
	//固有のセッションIDを作成
	sessionID := uuid.New().String()
	//sessionをDBに保存
	sessionDoc := bson.D{
		{"session_id",sessionID},
		{"user_id",userDocID},
	}
	_, err = dbhandler.Insert("googroutes", "sessions", sessionDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Println(err)
		return
	}

	signedStr,err := createToken(sessionID)
	if err != nil {
		http.Redirect(w,req,"/",http.StatusSeeOther)
		log.Println(err)
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