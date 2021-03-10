package tplutil

import (
	"app/controllers"
	"app/customerr"
	"app/model"
	"fmt"
	"net/http"
)

type TplData struct {
	Data map[string]interface{}
	User model.User
	Err  error
}

func getDataFromCtx(req *http.Request) *TplData {
	data, ok := req.Context().Value("data").(map[string]interface{})
	if !ok {
		return &TplData{
			Err: customerr.BaseErr{
				Op:  "Getting data from context",
				Msg: "エラーが発生しました。",
				Err: fmt.Errorf("error while getting data from context"),
			},
		}
	}
	return &TplData{
		Data: data,
	}
}

func GetTplData(req *http.Request) *TplData {
	tData := getDataFromCtx(req)
	tData.User, tData.Err = controllers.GetUserFromCtx(req)
	return tData
}
