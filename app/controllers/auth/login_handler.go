package auth

import (
	"app/bsonconv"
	"app/controllers"
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type loginController struct {
	controllers.Controller
	email    string
	password string
	user     model.User
}

//getUserFromDB fetches user document from DB
func (l *loginController) getUserFromDB() bson.M {
	if l.Err != nil {
		return nil
	}
	d, err := model.FindUser("email", l.email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			l.Err = customerr.BaseErr{
				Op:  "get user from DB",
				Msg: "メールアドレスまたはパスワードが間違っています。",
				Err: fmt.Errorf("error while getting user from DB: %w", err),
			}
		default:
			l.Err = customerr.BaseErr{
				Op:  "get user from DB",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while getting user from DB: %w", err),
			}
		}
		return nil
	}
	return d
}

//comparePassword compares hashed password and password user inputted
func (l *loginController) comparePasswords() {
	if l.Err != nil {
		return
	}
	err := bcrypt.CompareHashAndPassword(l.user.Password, []byte(l.password))
	if err != nil {
		l.Err = customerr.BaseErr{
			Op:  "compare passwords",
			Msg: "メールアドレスまたはパスワードが間違っています。",
			Err: fmt.Errorf("error while comparing passwords: %w", err),
		}
		return
	}
}

//generateNewSession generates new sessionID and save it to DB
func (l *loginController) generateNewSession(w http.ResponseWriter) {
	if l.Err != nil {
		return
	}
	err := genNewSession(l.user.ID, w)

	if err != nil {
		l.Err = customerr.BaseErr{
			Op:  "generate new session",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while generating new session: %w", err),
		}
		return
	}
}

func Login(w http.ResponseWriter, req *http.Request) {
	var l loginController
	l.GetStrValueFromCtx(req, &l.email, "email")
	l.GetStrValueFromCtx(req, &l.password, "password")
	d := l.getUserFromDB()
	bsonconv.DocToStruct(d, &l.user, &l.Err, "login user")
	l.comparePasswords()
	l.generateNewSession(w)

	if l.Err != nil {
		e := l.Err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/login_form")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
