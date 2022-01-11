package cookiehandler

import (
	"encoding/base64"
	"net/http"
	"regexp"
)

func MakeCookieAndRedirect(w http.ResponseWriter, req *http.Request, cName, cVal, path string) {
	b64CVal := base64.StdEncoding.EncodeToString([]byte(cVal))
	re := regexp.MustCompile(`(\w*/?)*`)
	//path omitted after query parameter
	pathWithoutQParam := re.FindString(path)
	c := &http.Cookie{
		Name:   cName,
		Value:  b64CVal,
		Path:   pathWithoutQParam,
		MaxAge: 1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, path, http.StatusSeeOther)
}

func DeleteCookie(w http.ResponseWriter, name, path string) {
	c := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   path,
		MaxAge: -1,
	}
	http.SetCookie(w, c)
}
