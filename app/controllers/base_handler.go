package controllers

import (
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

//CheckHTTPMethod checks HTTP method
func CheckHTTPMethod(req *http.Request, e *error) {
	if req.Method != "POST" {
		*e = customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method"),
		}
	}
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

//ConvertDucToStruct converts a bson document to a struct
func ConvertDucToStruct(d bson.M, s interface{}, e *error, modelName string) {
	if *e != nil {
		return
	}
	b, err := bson.Marshal(d)
	if err != nil {
		*e = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson marshaling %s: %w", modelName, err),
		}
	}
	err = bson.Unmarshal(b, s)
	if err != nil {
		*e = customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling %s: %w", modelName, err),
		}
	}
}

func GetLoginDataFromCtx(req *http.Request) map[string]interface{} {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		op := "Getting data from context"
		err := "error while getting data from context"
		log.Printf("operation: %s, error: %v", op, err)
		return map[string]interface{}{"isLoggedIn": false}
	}
	return data
}
