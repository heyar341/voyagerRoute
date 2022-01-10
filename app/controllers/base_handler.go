package controllers

import (
	"app/internal/customerr"
	"fmt"
	"net/http"
)

//CheckHTTPMethod checks HTTP method
func CheckHTTPMethod(req *http.Request, err *error) {
	if req.Method != "POST" {
		*err = customerr.BaseErr{
			Op:  "check HTTP method",
			Msg: "HTTPメソッドが不正です。",
			Err: fmt.Errorf("invalid HTTP method"),
		}
	}
}
