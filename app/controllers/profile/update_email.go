package profile

import (
	"app/cookiehandler"
	"app/customerr"
	"app/mailhandler"
	"app/model"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

const REDIRECT_URI_TO_UPDATE_EMAIL_FORM = "/profile/email_edit_form"

type updateEmailProcess struct {
	user     model.User
	newEmail string
	err      error
}

func (u *updateEmailProcess) checkHTTPMethod(req *http.Request) {
	if req.Method != "POST" {
		u.err = customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP access method"),
		}
	}
}

//Auth middlewareからuserIDを取得
func (u *updateEmailProcess) getUserFromCtx(req *http.Request) {
	if u.err != nil {
		return
	}
	user, ok := req.Context().Value("user").(model.User)
	if !ok {
		u.err = customerr.BaseErr{
			Op:  "get user from request's context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
		return
	}
	u.user = user
}

//リクエストフォームからメールアドレスを取得
func (u *updateEmailProcess) getEmailFromForm(req *http.Request) {
	if u.err != nil {
		return
	}
	newEmail := req.FormValue("email")
	if !mailhandler.IsEmailValid(newEmail) {
		u.err = customerr.BaseErr{
			Op:  "check email address's validity",
			Msg: "メールアドレスに不備があります。",
			Err: fmt.Errorf("request email was invalid %v", newEmail),
		}
		return
	}
	u.newEmail = newEmail
}

//emailを仮変更としてDBに保存
func (u *updateEmailProcess) saveEditingEmail(token string) {
	if u.err != nil {
		return
	}
	err := model.SaveEditingEmail(u.newEmail, token)
	if err != nil {
		u.err = customerr.BaseErr{
			Op:  "Saving editing email to DB",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting editing email to editing_email collection %w", err),
		}
		return
	}
}

func UpdateEmail(w http.ResponseWriter, req *http.Request) {
	var u updateEmailProcess
	u.checkHTTPMethod(req)
	u.getUserFromCtx(req)
	u.getEmailFromForm(req)
	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()
	u.saveEditingEmail(token)

	if u.err != nil {
		e := u.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_EMAIL_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	//認証依頼画面表示
	http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

	//メールでトークン付きのURLを送る
	err := mailhandler.SendConfirmEmail(token, u.newEmail, u.user.UserName, "confirm_email")
	return
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}

}

type confirmUpdateEmail struct {
	userID       primitive.ObjectID
	editingEmail model.EditingEmail
	token        string
	err          error
}

//Auth middlewareからuserIDを取得
func (c *confirmUpdateEmail) getUserFromCtx(req *http.Request) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok {
		log.Printf(": %v", ok)
		c.err = customerr.BaseErr{
			Op:  "get user from request's context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
	}
	c.userID = user.ID
}

//メール認証トークンをリクエストURLから取得
func (c *confirmUpdateEmail) getTokenFromURL(req *http.Request) {
	if c.err != nil {
		return
	}
	token := req.URL.Query()["token"][0]
	if token == "" {
		c.err = customerr.BaseErr{
			Op:  "get token form URL",
			Msg: "トークン情報が不正です。",
			Err: fmt.Errorf("token was empty"),
		}
	}
	c.token = token
}

func (c *confirmUpdateEmail) getEditingEmailDocFromDB() bson.M {
	if c.err != nil {
		return nil
	}
	d, err := model.GetEditingEmailDoc(c.token)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			c.err = customerr.BaseErr{
				Op:  "get editing email document form DB",
				Msg: "認証トークンが一致しません。",
				Err: fmt.Errorf("token was invalid: %w", err),
			}
		default:
			c.err = customerr.BaseErr{
				Op:  "get editing email document form DB",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while finding editing email document from editing_email collection %w", err),
			}
		}
		return nil
	}
	return d
}

func (c *confirmUpdateEmail) convertBSONToStruct(d bson.M) {
	if c.err != nil {
		return
	}
	b, err := bson.Marshal(d)
	if err != nil {
		c.err = customerr.BaseErr{
			Op:  "bson marshal editing_email document",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson marshaling editing_email document: %w", err),
		}
		return
	}
	err = bson.Unmarshal(b, &c.editingEmail)
	if err != nil {
		c.err = customerr.BaseErr{
			Op:  "bson unmarshal editing_email []byte document",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling []byte editing_email document: %w", err),
		}
		return
	}
}

//トークンの有効期限を確認
func (c *confirmUpdateEmail) checkTokenExpire() {
	if c.err != nil {
		return
	}
	t := time.Unix(c.editingEmail.ExpiresAt, 0)
	if !t.After(time.Now()) {
		c.err = customerr.BaseErr{
			Op:  "check if token expires or not",
			Msg: "トークンの有効期限が切れています。もう一度メールアドレス変更のお手続きをしてください。",
			Err: fmt.Errorf("token expired"),
		}
		return
	}
}

func (c *confirmUpdateEmail) updateUserEmail() {
	if c.err != nil {
		return
	}
	err := model.UpdateUser(c.userID, "email", c.editingEmail.Email)
	if err != nil {
		c.err = customerr.BaseErr{
			Op:  "update use document's email field",
			Msg: "メールアドレスの更新に失敗しました。",
			Err: fmt.Errorf("err while updating email field in user document: %w", err),
		}
		return
	}
}

func ConfirmUpdateEmail(w http.ResponseWriter, req *http.Request) {
	var c confirmUpdateEmail
	c.getUserFromCtx(req)
	c.getTokenFromURL(req)
	d := c.getEditingEmailDocFromDB()
	c.convertBSONToStruct(d)
	c.checkTokenExpire()
	c.updateUserEmail()

	if c.err != nil {
		e := c.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, REDIRECT_URI_TO_UPDATE_EMAIL_FORM)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	cookiehandler.MakeCookieAndRedirect(w, req, "success¬", "メールアドレスの変更が完了しました。", "/profile")
}
