package mailhandler

import (
	"app/envhandler"
	"app/model"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
)

func SendConfirmEmail(token, email, userName, path string) error {
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

	tokenURL := "グーグる〜とをご利用いただきありがとうございます。\n" +
		"このメールはメールアドレス認証用に送信されたメールです。\n" +
		"このメールを受信してから１時間以内に認証を行ってください。\n" +
		"１時間以内に認証が行われない場合、認証はキャンセルされます。\n\n" +
		"認証用URL:\n" +
		"https://googroutes.com/" + path + "/?token=" + token
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		mailAuth,
		"app.goog.routes@gmail.com",
		[]string{email},
		[]byte(fmt.Sprintf("To:%s\r\nSubject:メールアドレス認証のお願い\r\n\r\n%s", userName, tokenURL)),
	)
	if err != nil {
		log.Printf("Error sending email for confirm registering: %v", err)
		return err
	}
	return nil
}

func SendQuestion(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}

	//envファイルからGmailのアプリパスワード取得
	gmailPassword, err := envhandler.GetEnvVal("GMAIL_APP_PASS")
	if err != nil {
		msg := "送信中にエラーが発生しました。"
		http.Redirect(w, req, "/question_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting gmail app password form env file: %v", err)
		return
	}

	user, ok := req.Context().Value("user").(model.User)
	if !ok {
		msg := "送信中にエラーが発生しました。"
		http.Redirect(w, req, "/question_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting userName from context: %v", err)
		return
	}

	userName := user.UserName
	email := user.Email

	qText := req.FormValue("question")

	mailAuth := smtp.PlainAuth(
		"",
		"app.goog.routes@gmail.com",
		gmailPassword,
		"smtp.gmail.com",
	)

	t := "問い合わせ\n" +
		"ユーザー名：" + userName + "\n" +
		"メールアドレス:" + email + "\n" +
		"質問内容：" + qText + "\n"

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		mailAuth,
		"app.goog.routes@gmail.com",
		[]string{"app.goog.routes@gmail.com"},
		[]byte(fmt.Sprintf("To:%s\r\nSubject:問い合わせ\r\n\r\n%s", "自分", t)),
	)
	if err != nil {
		msg := "送信中にエラーが発生しました。"
		http.Redirect(w, req, "/question_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error sending email for question: %v", err)
		return
	}

	http.Redirect(w, req, "/mypage", http.StatusSeeOther)
}
