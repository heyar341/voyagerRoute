package auth

import (
	"app/cookiehandler"
	"app/customerr"
	"app/mailhandler"
	"app/model"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type registerProcess struct {
	userName        string
	email           string
	password        string
	securedPassword []byte
	err             error
}

func (r *registerProcess) getUserName(req *http.Request) {
	//Validation完了後のメールアドレスを取得
	username, ok := req.Context().Value("username").(string)
	if !ok {
		r.err = customerr.BaseErr{
			Op:  "get username from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting username from request context"),
		}
		return
	}
	r.userName = username
}

func (r *registerProcess) getEmail(req *http.Request) {
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		r.err = customerr.BaseErr{
			Op:  "get email from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting email from request context"),
		}
		return
	}
	r.email = email
}

func (r *registerProcess) getPassword(req *http.Request) {
	if r.err != nil {
		return
	}
	//Validation完了後のパスワードを取得
	password, ok := req.Context().Value("password").(string)
	if !ok {
		r.err = customerr.BaseErr{
			Op:  "get password from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting password from request context"),
		}
		return
	}

	r.password = password
}

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

//メール認証トークンをリクエストURLから取得
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

func (cR *confirmRegister) convertDucToStruct(d bson.M) {
	if cR.err != nil {
		return
	}
	b, err := bson.Marshal(d)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson marshaling registeringUser: %w", err),
		}
		return
	}
	err = bson.Unmarshal(b, &cR.registeringUser)
	if err != nil {
		cR.err = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling registeringUser: %w", err),
		}
		return
	}
}

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
	r.getUserName(req)
	r.getEmail(req)
	r.getPassword(req)
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
	err := mailhandler.SendConfirmEmail(token, r.email, r.userName)
	if err != nil {
		log.Printf("Error while sending email at registering: %v", err)
	}
}

func ConfirmRegister(w http.ResponseWriter, req *http.Request) {
	var cR confirmRegister
	cR.getToken(req)
	d := cR.findUserByToken()
	cR.convertDucToStruct(d)
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
