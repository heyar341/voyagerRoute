package profile

import (
	"app/controllers/envhandler"
	"app/dbhandler"
	"app/model"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
)

func UpdateEmail(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		msg = "リクエストメソッドが不正です。"
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}

	newEmail := req.FormValue("email")

	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()

	//emailを仮変更としてDBに保存
	//保存するドキュメント
	editingDoc := bson.D{
		{"email", newEmail},
		{"token", token},
	}
	//DBに保存
	_, err := dbhandler.Insert("googroutes", "editing_email", editingDoc)
	if err != nil {
		msg = "エラーが発生しました。しばらく経ってからもう一度ご利用ください。"
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg+"&email="+newEmail, http.StatusSeeOther)
		log.Println(err)
		return
	}
	//メール送信に少し時間がかかるので、認証依頼画面表示を先に処理
	http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

	//「メールでトークン付きのURLを送る」
	//envファイルからGmailのアプリパスワード取得
	gmailPassword, err := envhandler.GetEnvVal("GMAIL_APP_PASS")
	if err != nil {
		msg = "エラーが発生しました。しばらく経ってからもう一度ご利用ください。"
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg+"&email="+newEmail, http.StatusSeeOther)
		return
	}
	mailAuth := smtp.PlainAuth(
		"",
		"app.goog.routes@gmail.com",
		gmailPassword,
		"smtp.gmail.com",
	)

	tokenURL := "http://localhost:8080/confirm_email/?token=" + token //localhostは本番で変更
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		mailAuth,
		"app.goog.routes@gmail.com",
		[]string{newEmail},
		[]byte(fmt.Sprintf("To:%s\r\nSubject:メールアドレス認証のお願い\r\n\r\n%s", user.UserName, tokenURL)),
	)
	if err != nil {
		log.Printf("Error while sending email:%v", err)
	}
}

func ConfirmUpdateEmail(w http.ResponseWriter, req *http.Request) {
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
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
	resp, err := dbhandler.Find("googroutes", "editing_email", tokenDoc, nil)
	if err != nil {
		msg = url.QueryEscape("認証トークンが一致しません。")
		http.Redirect(w, req, "/?msg="+msg, http.StatusSeeOther)
		return
	}
	//DBから取得した値をmarshal
	bsonByte, _ := bson.Marshal(resp)

	type NewEmail struct {
		Email string `bson:"email"`
		Token string `bson:"token"`
	}

	var confirmedEmail NewEmail
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &confirmedEmail)
	//user documentを更新
	userDoc := bson.M{"_id": userID}
	updateDoc := bson.D{{"email", confirmedEmail.Email}}
	err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateDoc)
	if err != nil {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg+"&email="+confirmedEmail.Email, http.StatusSeeOther)
		log.Printf("Error while updating email: %v", err)
		return
	}

	success := "メールアドレスの変更が完了しました。"
	http.Redirect(w, req, "/?success="+success, http.StatusSeeOther)
}
