package mailhandler

import (
	"app/controllers/envhandler"
	"fmt"
	"log"
	"net/smtp"
)

func SendConfirmEmail(token, email, userName string) error {
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
		"https://googroutes.com/confirm_register/?token=" + token
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
