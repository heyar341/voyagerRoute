package mailhandler

import (
	"app/internal/cookiehandler"
	"app/internal/envhandler"
	"app/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var mailDomain = "mail.googroutes.com"

func SendConfirmEmail(token, email, path string) error {
	apiKey, err := envhandler.GetEnvVal("MAILGUN_API_KEY")
	if err != nil {
		return fmt.Errorf("couldn't get mailgun apiKey form env file: %w", err)
	}

	mg := mailgun.NewMailgun(mailDomain, apiKey)

	sender := "グーグる〜と運営 <customer_service@mail.googroutes.com>"
	subject := "メールアドレス認証のお願い。"
	body := "グーグる〜とをご利用いただきありがとうございます。\n" +
		"このメールはメールアドレス認証用に送信されたメールです。\n" +
		"このメールを受信してから１時間以内に認証を行ってください。\n" +
		"１時間以内に認証が行われない場合、認証はキャンセルされます。\n\n" +
		"認証用URL:\n" +
		"https://googroutes.com/" + path + "/?token=" + token
	recipient := email

	m := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = mg.Send(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func SendQuestion(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}

	user, ok := req.Context().Value("user").(model.User)
	if !ok {
		msg := "送信中にエラーが発生しました。"
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", msg, "/question_form")
		log.Printf("Error while getting userName from context")
		return
	}

	userName := user.UserName
	email := user.Email
	qText := req.FormValue("question")

	apiKey, err := envhandler.GetEnvVal("MAILGUN_API_KEY")
	if err != nil {
		msg := "送信中にエラーが発生しました。"
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", msg, "/question_form")
		log.Printf("couldn't get mailgun apiKey form env file: %v", err)
		return
	}
	myMail, err := envhandler.GetEnvVal("MY_MAIL")
	if err != nil {
		msg := "送信中にエラーが発生しました。"
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", msg, "/question_form")
		log.Printf("couldn't get mailgun apiKey form env file: %v", err)
		return
	}

	mg := mailgun.NewMailgun(mailDomain, apiKey)

	sender := "グーグる〜と運営 <customer_service@mail.googroutes.com>"
	subject := "お問い合わせ"
	body := "お問い合わせ\n" +
		"ユーザー名：" + userName + "\n" +
		"メールアドレス:" + email + "\n" +
		"質問内容：" + qText + "\n"
	recipient := myMail

	m := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err = mg.Send(ctx, m)
	if err != nil {
		msg := "送信中にエラーが発生しました。"
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", msg, "/question_form")
		log.Printf("Error sending email for question: %v", err)
		return
	}

	cookiehandler.MakeCookieAndRedirect(w, req, "success", "お問い合わせの受付が完了しました。", "/mypage")
}
