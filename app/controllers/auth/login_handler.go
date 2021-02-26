package auth

import (
	"app/dbhandler"
	"app/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/url"
)

func Login(w http.ResponseWriter, req *http.Request) {
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}
	password, ok := req.Context().Value("password").(string)
	//Validation完了後のパスワードを取得
	if !ok {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
	}

	//取得するドキュメントの条件
	emailDoc := bson.D{{"email", email}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "users", emailDoc, nil)
	if err != nil {
		msg := "メールアドレスまたはパスワードが正しくありません。"
		//入力されたメールアドレスを保持する
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		return
	}
	//DBから取得した値をmarshal
	bsonByte, _ := bson.Marshal(resp)
	var user model.UserData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	storedPass := user.Password
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(password))
	//一致しない場合
	if err != nil {
		msg := "メールアドレスまたはパスワードが正しくありません。"
		//入力されたメールアドレスを保持する
		email = url.QueryEscape(email)
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		return
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
		http.Redirect(w, req, "/login_form/?msg="+msg, http.StatusSeeOther)
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
		Path:  "/",
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, "/", http.StatusSeeOther)
}
