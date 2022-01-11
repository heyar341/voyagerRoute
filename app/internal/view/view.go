package view

import (
	"encoding/base64"
	"html/template"
	"net/http"
)

func ShowMsgWithCookie(w http.ResponseWriter, c *http.Cookie, data map[string]interface{}, tType *template.Template, tName string) {
	b64Str, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		tType.ExecuteTemplate(w, tName, data)
		return
	}
	data[c.Name] = string(b64Str)
	tType.ExecuteTemplate(w, tName, data)
}

func ExistsCookie(w http.ResponseWriter, req *http.Request, data map[string]interface{}, tType *template.Template, tName string) bool {
	//successメッセージがある場合
	c, _ := req.Cookie("success")
	if c != nil {
		ShowMsgWithCookie(w, c, data, tType, tName)
		return true
	}
	//エラーメッセージがある場合
	c, _ = req.Cookie("msg")
	if c != nil {
		ShowMsgWithCookie(w, c, data, tType, tName)
		return true
	}
	return false
}
