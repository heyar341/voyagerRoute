package auth

import (
	"app/cookiehandler"
	"app/customerr"
	"app/model"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type loginProcess struct {
	email    string
	password string
	user     model.User
	err      error
}

//getEmail gets email from request form
func (l *loginProcess) getEmail(req *http.Request) {
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		l.err = customerr.BaseErr{
			Op:  "get email from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting email from request context"),
		}
		return
	}
	l.email = email
}

//getPassword gets password from request form
func (l *loginProcess) getPassword(req *http.Request) {
	if l.err != nil {
		return
	}
	//Validation完了後のパスワードを取得
	password, ok := req.Context().Value("password").(string)
	if !ok {
		l.err = customerr.BaseErr{
			Op:  "get password from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting password from request context"),
		}
		return
	}

	l.password = password
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

//convertDocToStruct converts user document to User struct
func (l *loginProcess) convertDocToStruct(d bson.M) {
	if l.err != nil {
		return
	}
	b, err := bson.Marshal(d)
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson marshaling user: %w", err),
		}
		return
	}
	err = bson.Unmarshal(b, &l.user)
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling user: %w", err),
		}
		return
	}
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
	l.getEmail(req)
	l.getPassword(req)
	d := l.getUserFromDB()
	l.convertDocToStruct(d)
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
