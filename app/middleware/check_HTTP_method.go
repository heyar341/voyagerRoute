package middleware

import (
	"log"
	"net/http"
)

//CheckHTTPMethod checks HTTP method
func CheckHTTPMethod(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "HTTPメソッドが不正です。", http.StatusMethodNotAllowed)
			log.Printf("Invalid HTTP request came")
			return
		}

		next.ServeHTTP(w, req)
	}
}

//CheckHTTPContentType checks request's Content-Type
func CheckHTTPContentType(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "HTTPメソッドが不正です。", http.StatusMethodNotAllowed)
			log.Printf("Invalid HTTP request came")
			return
		}
		next.ServeHTTP(w, req)
	}
}
