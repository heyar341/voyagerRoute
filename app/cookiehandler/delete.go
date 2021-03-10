package cookiehandler

import "net/http"

func DeleteCookie(w http.ResponseWriter, name, path string) {
	c := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   path,
		MaxAge: -1,
	}
	http.SetCookie(w, c)
}
