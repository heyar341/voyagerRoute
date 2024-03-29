package auth

import (
	"app/controllers"
	"app/internal/bsonconv"
	"app/internal/cookiehandler"
	"app/internal/customerr"
	"app/internal/errormsg"
	"app/internal/mailhandler"
	"app/internal/view"
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

type registerController struct {
	controllers.Controller
	userName        string
	email           string
	password        string
	securedPassword []byte
}

//generateSecuredPassword generates a hashed password
func (r *registerController) generateSecuredPassword() {
	if r.Err != nil {
		return
	}
	//パスワードをハッシュ化
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(r.password), 5)
	if err != nil {
		r.Err = customerr.BaseErr{
			Op:  "generate hashed password",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while hashing password :%w", err),
		}
		return
	}
	r.securedPassword = securedPassword
}

//saveRegisteringUserToDB inserts user document to DB
func (r *registerController) saveRegisteringUserToDB(token string) {
	if r.Err != nil {
		return
	}
	//DBに保存
	err := model.SaveRegisteringUser(r.userName, r.email, token, r.securedPassword)
	if err != nil {
		r.Err = customerr.BaseErr{
			Op:  "insert user to registering collection",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while inserting user to registering collecion :%w", err),
		}
		return
	}
}

type confirmRegisterController struct {
	registeringUser model.Registering
	token           string
	userID          primitive.ObjectID
	err             error
}

//getTokenFromURL gets token from query parameter
func (cR *confirmRegisterController) getTokenFromURL(req *http.Request) {
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
func (cR *confirmRegisterController) findUserByToken() bson.M {
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
func (cR *confirmRegisterController) checkTokenExpire() {
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
func (cR *confirmRegisterController) saveNewUserToDB() {
	if cR.err != nil {
		return
	}
	userID, err := model.SaveNewUser(cR.registeringUser.UserName, cR.registeringUser.Email, cR.registeringUser.Password)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "insert user to users collection",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while inserting user to users collection: %w", err),
		}
		return
	}
	cR.userID = userID
}

//generateNewSession generates new session
func (cR *confirmRegisterController) generateNewSession(w http.ResponseWriter) {
	if cR.err != nil {
		return
	}
	err := genNewSession(cR.userID, w)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "generate new session",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while generating a new session: %w", err),
		}
		return
	}
}

func Register(w http.ResponseWriter, req *http.Request) {
	var r registerController
	if req.Method == "GET" {
		data := map[string]interface{}{"isLoggedIn": false}
		c, err := req.Cookie("msg")
		if err == nil {
			view.ShowMsgWithCookie(w, c, data, authTpl, "register.html")
			return
		}
		authTpl.ExecuteTemplate(w, "register.html", data)
		return
	}
	r.GetStrValueFromCtx(req, &r.userName, "username")
	r.GetStrValueFromCtx(req, &r.email, "email")
	r.GetStrValueFromCtx(req, &r.password, "password")
	r.generateSecuredPassword()
	//メールアドレス認証用のトークンを作成
	token := uuid.New().String()
	r.saveRegisteringUserToDB(token)
	if r.Err != nil {
		e := r.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/register_form")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	//認証依頼画面表示
	data := map[string]interface{}{"isLoggedIn": false} //新規登録はログイン時にしないと想定
	authTpl.ExecuteTemplate(w, "ask_confirm_email.html", data)

	//「メールでトークン付きのURLを送る」
	err := mailhandler.SendConfirmEmail(token, r.email, "confirm_register")
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}
}

func ConfirmRegister(w http.ResponseWriter, req *http.Request) {
	var cR confirmRegisterController
	cR.getTokenFromURL(req)
	d := cR.findUserByToken()
	bsonconv.DocToStruct(d, &cR.registeringUser, &cR.err, "registeringUser")
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
