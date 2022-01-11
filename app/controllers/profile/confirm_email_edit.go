package profile

import (
	"app/controllers"
	"app/internal/bsonconv"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type confirmUpdateEmailController struct {
	controllers.Controller
	user         model.User
	editingEmail model.EditingEmail
	token        string
	err          error
}

//getTokenFromURL gets token from URL parameter
func (c *confirmUpdateEmailController) getTokenFromURL(req *http.Request) {
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

//getEditingEmailDocFromDB fetch editing email from DB
func (c *confirmUpdateEmailController) getEditingEmailDocFromDB() bson.M {
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

//checkTokenExpire checks if token expires or not
func (c *confirmUpdateEmailController) checkTokenExpire() {
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

//updateUserEmail updates email field in user document
func (c *confirmUpdateEmailController) updateUserEmail() {
	if c.err != nil {
		return
	}
	err := model.UpdateUser(c.user.ID, "email", c.editingEmail.Email)
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
	var c confirmUpdateEmailController
	c.GetUserFromCtx(req, &c.user)
	c.getTokenFromURL(req)
	d := c.getEditingEmailDocFromDB()
	bsonconv.DocToStruct(d, &c.editingEmail, &c.err, "editing email")
	c.checkTokenExpire()
	c.updateUserEmail()

	if c.err != nil {
		e := c.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, redirectURIToUpdateEmailForm)
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}
	cookiehandler.MakeCookieAndRedirect(w, req, "success¬", "メールアドレスの変更が完了しました。", "/profile")
}
