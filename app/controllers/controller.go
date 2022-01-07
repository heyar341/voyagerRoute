package controllers

import (
	"app/customerr"
	"app/model"
	"fmt"
	"log"
	"net/http"
)

type Controller struct {
	Err error
}

type ControllerInterface interface {
	GetLoginStateFromCtx(req *http.Request)
	GetUserFromCtx(req *http.Request, user *model.User)
	GetStrValueFromCtx(req *http.Request, field *string, valueName string)
}

func (c *Controller) GetLoginStateFromCtx(req *http.Request) map[string]interface{} {
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
func (c *Controller) GetUserFromCtx(req *http.Request, user *model.User) {
	if c.Err != nil {
		return
	}
	u, ok := req.Context().Value("user").(model.User)
	if !ok {
		c.Err = customerr.BaseErr{
			Op:  "get user from request's context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting user from reuest's context"),
		}
		return
	}
	*user = u
}

//GetStrValueFromCtx gets a string value from request's context
func (c *Controller) GetStrValueFromCtx(req *http.Request, field *string, valueName string) {
	if c.Err != nil {
		return
	}
	//Validation完了後の値を取得
	v, ok := req.Context().Value(valueName).(string)
	if !ok {
		c.Err = customerr.BaseErr{
			Op:  "get" + valueName + "from request context",
			Msg: "エラーが発生しました。",
			Err: fmt.Errorf("error while getting %s from request context", valueName),
		}
		return
	}
	*field = v
}
