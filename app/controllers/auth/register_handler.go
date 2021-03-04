package auth

import (
	"app/dbhandler"
	"app/mailhandler"
	"encoding/json"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
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
		{"expires_at", time.Now().Add(1 * time.Hour).Unix()},
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
	err = mailhandler.SendConfirmEmail(token, email, userName)
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}
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

	expireBSON, ok := userM["expires_at"].(primitive.DateTime)
	if !ok {
		msg := url.QueryEscape("データの処理中にエラーが発生しました。申し訳ありませんが、もう一度新規登録操作をお願いいたします。")
		http.Redirect(w, req, "/profile/register_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while confirming register. type asserting token's expires time: %v", err)
		return
	}
	expires := expireBSON.Time() //time.Time

	//トークンの有効期限を確認
	if expires.After(time.Now()) {
		msg := url.QueryEscape("トークンの有効期限が切れています。もう一度新規登録操作をお願いいたします。")
		http.Redirect(w, req, "/profile/register_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("registering token of %v is expired", userM["email"])
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
