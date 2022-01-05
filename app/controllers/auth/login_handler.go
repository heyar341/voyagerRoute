package auth

import (
	"app/bsonconv"
	"app/contexthandler"
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

type loginProcess struct {
	email    string
	password string
	user     model.User
	err      error
}

//getUserFromDB fetches user document from DB
func (l *loginProcess) getUserFromDB() bson.M {
	if l.err != nil {
		return nil
	}
	d, err := model.FindUser("email", l.email)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			l.err = customerr.BaseErr{
				Op:  "get user from DB",
				Msg: "メールアドレスまたはパスワードが間違っています。",
				Err: fmt.Errorf("error while getting user from DB: %w", err),
			}
		default:
			l.err = customerr.BaseErr{
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
func (l *loginProcess) comparePasswords() {
	if l.err != nil {
		return
	}
	err := bcrypt.CompareHashAndPassword(l.user.Password, []byte(l.password))
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "compare passwords",
			Msg: "メールアドレスまたはパスワードが間違っています。",
			Err: fmt.Errorf("error while comparing passwords: %w", err),
		}
		return
	}
}

//generateNewSession generates new sessionID and save it to DB
func (l *loginProcess) generateNewSession(w http.ResponseWriter) {
	if l.err != nil {
		return
	}
	err := genNewSession(l.user.ID, w)

	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "generate new session",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while generating new session: %w", err),
		}
		return
	}
}

func Login(w http.ResponseWriter, req *http.Request) {
	var l loginProcess
	contexthandler.GetStrValueFromCtx(req, &l.email, &l.err, "email")
	contexthandler.GetStrValueFromCtx(req, &l.password, &l.err, "password")
	d := l.getUserFromDB()
	bsonconv.ConvertDucToStruct(d, &l.user, &l.err, "login user")
	l.comparePasswords()
	l.generateNewSession(w)

	if l.err != nil {
		e := l.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/login_form")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
