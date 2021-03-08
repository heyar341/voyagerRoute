package tplutil

import (
	"fmt"
	"app/model"
	"app/customerr"
	"net/http"
)

type TplData struct {
	Data map[string]interface{}
	User model.UserData
	Err error
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

func (t *TplData) getUserFromCtx(req *http.Request) {
	if t.Err != nil {
		return
	}
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		t.Err = customerr.BaseErr{
			Op:  "Getting user from context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from context"),
		}
		return
	}
	t.User = user
}

func GetTplData(req *http.Request) *TplData {
	tData := getDataFromCtx(req)
	tData.getUserFromCtx(req)
	return tData
}
