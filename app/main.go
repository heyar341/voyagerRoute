// Copyright 2015 Google Inc. All Rights Reserved.
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"html/template"
	"app/controllers/auth"
	"app/controllers/routes"

)
var tpl *template.Template

var auth_tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/route_search/*"))
	auth_tpl = template.Must(template.ParseGlob("templates/auth/*"))
}

func main() {

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	//Authentication
	http.HandleFunc("/register_form",registerForm)
	http.HandleFunc("/register",auth.Register)
	http.HandleFunc("/login_form",loginForm)
	http.HandleFunc("/login",auth.Login)
	//Direction API
	http.HandleFunc("/show_map",index)
	http.HandleFunc("/routes_save",routes.SaveRoutes)

	http.ListenAndServe(":80",nil)
}

func registerForm(w http.ResponseWriter, req *http.Request) {
	auth_tpl.ExecuteTemplate(w, "register.html",nil)
}
func loginForm(w http.ResponseWriter, req *http.Request) {
	auth_tpl.ExecuteTemplate(w, "login.html",nil)
}
func index(w http.ResponseWriter, req *http.Request){
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil{
		panic("Can't load env file")
	}
	//envファイルからAPI key取得
	apiKey := os.Getenv("MAP_API_KEY")
	data := map[string]string{"apiKey":apiKey}
	tpl.ExecuteTemplate(w, "place_and_direction_improve.html", data)
}

