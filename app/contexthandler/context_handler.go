package contexthandler

import (
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

func GetLoginStateFromCtx(req *http.Request) map[string]interface{} {
	isLoggedIn, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		op := "Getting data from context"
		err := "error while getting data from context"
		log.Printf("operation: %s, error: %v", op, err)
		return map[string]interface{}{"isLoggedIn": false}
	}
	return isLoggedIn
}

//GetUserFromCtx gets user from Auth middleware
func GetUserFromCtx(req *http.Request, user *model.User, e *error) {
	if *e != nil {
		return
	}
	u, ok := req.Context().Value("user").(model.User)
	if !ok {
		*e = customerr.BaseErr{
			Op:  "get user from request's context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
		return
	}
	*user = u
}

//GetStrValueFromCtx gets a string value from request's context
func GetStrValueFromCtx(req *http.Request, field *string, e *error, valueName string) {
	if *e != nil {
		return
	}
	//Validation完了後の値を取得
	v, ok := req.Context().Value(valueName).(string)
	if !ok {
		*e = customerr.BaseErr{
			Op:  "get" + valueName + "from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting %s from request context", valueName),
		}
		return
	}
	*field = v
}
