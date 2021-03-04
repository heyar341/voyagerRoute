package profile

import (
	"app/dbhandler"
	"app/mailhandler"
	"app/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"net/url"
	"time"
)

func UpdateEmail(w http.ResponseWriter, req *http.Request) {
	msg := "エラーが発生しました。もう一度操作を行ってください。"
	if req.Method != "POST" {
		msg = "リクエストメソッドが不正です。"
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}

	newEmail := req.FormValue("email")
	if !mailhandler.IsEmailValid(newEmail) {
		msg = url.QueryEscape("メールアドレスに不備があります。")
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg+"&newEmail="+url.QueryEscape(newEmail), http.StatusSeeOther)
		return
	}

	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()

	//emailを仮変更としてDBに保存
	//保存するドキュメント
	editingDoc := bson.D{
		{"email", newEmail},
		{"expires_at", time.Now().Add(1 * time.Hour).Unix()},
		{"token", token},
	}
	//editing_email collectionに保存
	_, err := dbhandler.Insert("googroutes", "editing_email", editingDoc)
	if err != nil {
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg+"&email="+url.QueryEscape(newEmail), http.StatusSeeOther)
		log.Println(err)
		return
	}
	//メール送信に少し時間がかかるので、認証依頼画面表示を先に処理
	http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

	//「メールでトークン付きのURLを送る」
	mailhandler.SendConfirmEmail(token, newEmail, user.UserName)
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}

}

func ConfirmUpdateEmail(w http.ResponseWriter, req *http.Request) {
	msg := url.QueryEscape("エラーが発生しました。もう一度操作を行ってください。")
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}
	userID := user.ID

	//メール認証トークンをリクエストURLから取得
	token := req.URL.Query()["token"][0]
	if token == "" {
		msg = "トークン情報が不正です。"
		http.Redirect(w, req, "/?msg=", http.StatusSeeOther)
		return
	}

	//このtokenはメール認証用でjwtを使ってないからParseTokenは呼び出さなくていい

	//取得するドキュメントの条件
	tokenDoc := bson.D{{"token", token}}
	//DBから取得
	em, err := dbhandler.Find("googroutes", "editing_email", tokenDoc, nil)
	if err != nil {
		msg = url.QueryEscape("認証トークンが一致しません。")
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		return
	}

	newEmail := em["email"]
	expireBSON, ok := em["expires_at"].(primitive.DateTime)
	if !ok {
		msg = url.QueryEscape("データの処理中にエラーが発生しました。申し訳ありませんが、もう一度メールアドレス変更のお手続きをしてください。")
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while type asserting token's expires time: %v", err)
		return
	}
	expires := expireBSON.Time() //time.Time

	//トークンの有効期限を確認
	if expires.After(time.Now()) {
		msg = url.QueryEscape("トークンの有効期限が切れています。もう一度メールアドレス変更のお手続きをしてください。")
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Editing_email token of %v is expired", newEmail)
		return
	}
	//user documentを更新
	userDoc := bson.M{"_id": userID}
	updateDoc := bson.D{{"email", newEmail}}
	err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateDoc)
	if err != nil {
		msg = url.QueryEscape("データの処理中にエラーが発生しました。申し訳ありませんが、もう一度メールアドレス変更のお手続きをしてください。")
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while updating email: %v", err)
		return
	}

	success := url.QueryEscape("メールアドレスの変更が完了しました。")
	http.Redirect(w, req, "/profile/?success="+success, http.StatusSeeOther)
}
