package auth

import (
	"app/controllers/envhandler"
	"app/dbhandler"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"time"
)

func Register(w http.ResponseWriter, req *http.Request) {
	//ユーザーに表示するエラーメッセージを定義
	msg := url.QueryEscape("エラ〜が発生しました。もう一度操作をしなおしてください。")
	//Validation完了後のユーザー名を取得
	userName, ok := req.Context().Value("username").(string)
	if !ok {
		http.Redirect(w, req, "/login_form/?msg="+msg+"&username="+userName, http.StatusSeeOther)
		return
	}
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		return
	}
	password, ok := req.Context().Value("password").(string)
	//Validation完了後のパスワードを取得
	if !ok {
		http.Redirect(w, req, "/login_form/?msg="+msg+"&email="+email, http.StatusSeeOther)
		return
	}

	//パスワードをハッシュ化
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+userName+"&email="+email, http.StatusSeeOther)
		log.Println(err)
		return
	}
	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()

	//userを仮登録としてDBに保存
	//保存するドキュメント
	registerDoc := bson.D{
		{"username", userName},
		{"email", email},
		{"password", securedPassword},
		{"token", token},
	}
	//DBに保存
	_, err = dbhandler.Insert("googroutes", "registering", registerDoc)
	if err != nil {
		http.Redirect(w, req, "/register_form/?msg="+msg+"&username="+userName+"&email="+email, http.StatusSeeOther)
		log.Printf("Error while inserting registering user to registering collection :%v", err)
		return
	}

	//メール送信に少し時間がかかるので、認証依頼画面表示を先に処理
	http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

	//「メールでトークン付きのURLを送る」
	err = sendConfirmEmail(token, email, userName)
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusInternalServerError)
		return
	}
}

func sendConfirmEmail(token, email, userName string) error {
	//envファイルからGmailのアプリパスワード取得
	gmailPassword, err := envhandler.GetEnvVal("GMAIL_APP_PASS")
	if err != nil {
		log.Printf("Error while getting gmail app password form env file: %v", err)
		return err
	}
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
		log.Printf("Error sending email for confir registering: %v", err)
		return err
	}
	return nil
}

func ConfirmRegister(w http.ResponseWriter, req *http.Request) {
	//メール認証トークンをリクエストURLから取得
	token := req.URL.Query()["token"][0]
	if token == "" {
		http.Redirect(w, req, "/register_form", http.StatusSeeOther)
		return
	}
	//このtokenはメール認証用でjwtを使ってないからParseTokenは呼び出さなくていい

	//取得するドキュメントの条件
	tokenDoc := bson.D{{"token", token}}

	//DBから取得
	userM, err := dbhandler.Find("googroutes", "registering", tokenDoc, nil)
	if err != nil {
		msg := url.QueryEscape("認証トークンが一致しません。")
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	//multi_route_titlesフィールドを追加
	userDoc := bson.D{{"username", userM["username"]},
		{"email", userM["email"]},
		{"password", userM["password"]},
		{"multi_route_titles", map[string]time.Time{}},
	}

	//「userをDBに保存」
	//DBに保存
	msg := url.QueryEscape("エラーが発生しました。もう一度操作をしなおしてください。")
	userID, err := dbhandler.Insert("googroutes", "users", userDoc)
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while inserting user data into DB: %v", err)
		return
	}

	//「session作成」
	err = genNewSession(userID, w)
	if err != nil {
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	succcess := url.QueryEscape("メールアドレス認証が完了しました。")
	http.Redirect(w, req, "/?success="+succcess, http.StatusSeeOther)
}

func EmailIsAvailable(w http.ResponseWriter, req *http.Request) {
	//メールアドレスが使用可能かのリクエスト
	type ValidEmailRequest struct {
		Email string `json:"email"`
	}

	if req.Method != "POST" {
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}
	//requestのフィールドを保存する変数
	var reqFields ValidEmailRequest
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &reqFields)
	if err != nil {
		http.Error(w, "入力に不正があります。", http.StatusBadRequest)
		log.Printf("Error while json marshaling: %v", err)
		return
	}

	//メールアドレスが使用可能か入れる変数
	var isValid = false
	emailDoc := bson.D{{"email", reqFields.Email}}
	//DBから取得
	_, err = dbhandler.Find("googroutes", "users", emailDoc, nil)
	//ドキュメントがない場合、メールアドレスは使用可能
	if err == mongo.ErrNoDocuments {
		isValid = true
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	type ResponseMsg struct {
		Valid bool `json:"valid"`
	}
	msg := ResponseMsg{Valid: isValid}
	respJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)

}
