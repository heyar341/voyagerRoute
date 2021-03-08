package cookiehandler

import (
	"encoding/base64"
	"net/http"
)

func MakeCookieAndRedirect(w http.ResponseWriter, req *http.Request, cName, cVal, path string) {
	b64CVal := base64.StdEncoding.EncodeToString([]byte(cVal))
	c := &http.Cookie{
		Name:   cName,
		Value:  b64CVal,
		Path:   path,
		MaxAge: 5,
	}
	http.SetCookie(w, c)
	http.Redirect(w, req, path, http.StatusSeeOther)
}
