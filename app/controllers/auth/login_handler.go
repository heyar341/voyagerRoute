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
	user     model.UserData
	err      error
}

func getEmail(req *http.Request) *loginProcess {
	//Validation完了後のメールアドレスを取得
	email, ok := req.Context().Value("email").(string)
	if !ok {
		return &loginProcess{
			err: customerr.BaseErr{
				Op:  "get email from request context",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while getting email from request context"),
			},
		}
	}
	return &loginProcess{
		email: email,
	}
}

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
	}

	l.password = password
}

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
	}
	err = bson.Unmarshal(b, &l.user)
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling user: %w", err),
		}
	}

}

func (l *loginProcess) comparePasswords() {
	if l.err != nil {
		return
	}
	err := bcrypt.CompareHashAndPassword(l.user.Password, []byte(l.password))
	if err != nil {
		l.err = customerr.BaseErr{
			Op:  "compare passwords",
			Msg: "メールアドレスまたはパスワードが間違っています。",
			Err: fmt.Errorf("error while bson marshaling user: %w", err),
		}
	}
}

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
	}
}

func Login(w http.ResponseWriter, req *http.Request) {
	l := getEmail(req)
	l.getPassword(req)
	//user documentを取得
	d := l.getUserFromDB()
	l.convertDocToStruct(d)
	l.comparePasswords()
	l.generateNewSession(w)

	if l.err != nil {
		e := l.err.(customerr.BaseErr)
		cookiehandler.MakeCookieAndRedirect(w, req, "msg", e.Msg, "/")
		log.Printf("operation: %s, error: %v", e.Op, e.Err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}