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
func CheckHTTPMethod(req *http.Request) error {
	if req.Method != "POST" {
		return customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method"),
		}
	}
	return nil
}

//GetStrValueFromCtx gets a string value from request's context
func GetStrValueFromCtx(req *http.Request, valueName string) (string, error) {
	//Validation完了後の値を取得
	v, ok := req.Context().Value(valueName).(string)
	if !ok {
		return "", customerr.BaseErr{
			Op:  "get" + valueName + "from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting %s from request context", valueName),
		}
	}
	return v, nil
}

//GetUserFromCtx gets user from Auth middleware
func GetUserFromCtx(req *http.Request) (model.User, error) {
	user, ok := req.Context().Value("user").(model.User)
	if !ok {
		return model.User{}, customerr.BaseErr{
			Op:  "get user from request's context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
	}
	return user, nil
}

//ConvertDucToStruct converts a bson document to a struct
func ConvertDucToStruct(d bson.M, s interface{}, modelName string) error {
	b, err := bson.Marshal(d)
	if err != nil {
		return customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson marshaling %s: %w", modelName, err),
		}
	}
	err = bson.Unmarshal(b, s)
	if err != nil {
		return customerr.BaseErr{
			Op:  "convert BSON document to struct",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while bson unmarshaling %s: %w", modelName, err),
		}
	}
	return nil
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
