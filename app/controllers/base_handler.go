package controllers

import (
	"app/customerr"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
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
