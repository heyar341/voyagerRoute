package auth

import (
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/mailhandler"
	"app/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type registerProcess struct {
	userName        string
	email           string
	password        string
	securedPassword []byte
	err             error
}

//generateSecuredPassword generates a hashed password
func (r *registerProcess) generateSecuredPassword() {
	if r.err != nil {
		return
	}
	//パスワードをハッシュ化
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(r.password), 5)
	if err != nil {
		r.err = customerr.BaseErr{
			Op:  "generate hashed password",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while hashing password :%w", err),
		}
		return
	}
	r.securedPassword = securedPassword
}

//saveRegisteringUser inserts user document to DB
func (r *registerProcess) saveRegisteringUser(token string) {
	if r.err != nil {
		return
	}
	//DBに保存
	err := model.SaveRegisteringUser(r.userName, r.email, token, r.securedPassword)
	if err != nil {
		r.err = customerr.BaseErr{
			Op:  "insert user to registering collection",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting user to registering collecion :%w", err),
		}
		return
	}
}

type confirmRegister struct {
	registeringUser model.Registering
	token           string
	userID          primitive.ObjectID
	err             error
}

//getToken gets token from query parameter
func (cR *confirmRegister) getToken(req *http.Request) {
	token := req.URL.Query()["token"][0]
	if token == "" {
		cR.err = customerr.BaseErr{
			Op:  "get token from url",
			Msg: "認証トークンがありません。",
			Err: fmt.Errorf("error while getting token from URL parameter"),
		}
		return
	}
	cR.token = token
}

//findUserByToken fetch user document from DB using token
func (cR *confirmRegister) findUserByToken() bson.M {
	if cR.err != nil {
		return nil
	}
	d, err := model.FindUserByToken(cR.token)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "find user by using token",
			Msg: "認証トークンが一致しません。",
			Err: fmt.Errorf("error while finding user from registering collecion :%w", err),
		}
		return nil
	}
	return d
}

//checkTokenExpire checks if token expires or not
func (cR *confirmRegister) checkTokenExpire() {
	if cR.err != nil {
		return
	}
	var t time.Time
	t = time.Unix(cR.registeringUser.ExpiresAt, 0)
	if !t.After(time.Now()) {
		cR.err = customerr.BaseErr{
			Op:  "check if token expires or not",
			Msg: "トークンの有効期限が切れています。もう一度新規登録操作をお願いいたします。",
			Err: fmt.Errorf("[%s]'s token expired", cR.registeringUser.Email),
		}
		return
	}
}

//saveNewUserToDB saves user document to users collection
func (cR *confirmRegister) saveNewUserToDB() {
	if cR.err != nil {
		return
	}
	userID, err := model.SaveNewUser(cR.registeringUser.UserName, cR.registeringUser.Email, cR.registeringUser.Password)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "insert user to users collection",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while inserting user to users collection: %w", err),
		}
		return
	}
	cR.userID = userID
}

//generateNewSession generates new session
func (cR *confirmRegister) generateNewSession(w http.ResponseWriter) {
	if cR.err != nil {
		return
	}
	err := genNewSession(cR.userID, w)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "generate new session",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while generating a new session: %w", err),
		}
		return
	}
}

func Register(w http.ResponseWriter, req *http.Request) {
	var r registerProcess
	r.userName, r.err = controllers.GetStrValueFromCtx(req, "username")
	r.email, r.err = controllers.GetStrValueFromCtx(req, "email")
	r.password, r.err = controllers.GetStrValueFromCtx(req, "password")
	r.generateSecuredPassword()
	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()
	r.saveRegisteringUser(token)
	if r.err != nil {
		e := r.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/register_form")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	//認証依頼画面表示
	http.Redirect(w, req, "/ask_confirm", http.StatusSeeOther)

	//「メールでトークン付きのURLを送る」
	err := mailhandler.SendConfirmEmail(token, r.email, r.userName, "confirm_register")
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}
}

func ConfirmRegister(w http.ResponseWriter, req *http.Request) {
	var cR confirmRegister
	cR.getToken(req)
	d := cR.findUserByToken()
	cR.err = controllers.ConvertDucToStruct(d, &cR.registeringUser, "registeringUser")
	cR.checkTokenExpire()
	cR.saveNewUserToDB()
	cR.generateNewSession(w)
	if cR.err != nil {
		e := cR.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	cookiehandler.MakeCookieAndRedirect(w, req, "success", "メールアドレス認証が完了しました。", "/")
}
