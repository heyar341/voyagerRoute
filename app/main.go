// Copyright 2015 Google Inc. All Rights Reserved.
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"html/template"
	"app/controllers/auth"
	"app/controllers/routes"

)
var tpl, auth_tpl, home_tpl, simul_search_tpl *template.Template

func init() {
	tpl = template.Must(template.Must(template.ParseGlob("templates/route_search/*")).ParseGlob("templates/includes/*.html"))
	simul_search_tpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/*")).ParseGlob("templates/includes/*.html"))
	auth_tpl = template.Must(template.Must(template.ParseGlob("templates/auth/*")).ParseGlob("templates/includes/*.html"))
	home_tpl = template.Must(template.Must(template.ParseGlob("templates/home/home.html")).ParseGlob("templates/includes/*.html"))
}

func main() {

	fmt.Println("App started")
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	//Authentication
	http.HandleFunc("/register_form/",registerForm)
	http.HandleFunc("/register",auth.Register)
	http.HandleFunc("/login_form/",loginForm)
	http.HandleFunc("/login",auth.Login)
	http.HandleFunc("/confirm_register/",auth.ConfirmRegister)
	http.HandleFunc("/ask_confirm/",askConfirm)

	//Direction API
	http.HandleFunc("/show_map",index)
	http.HandleFunc("/routes_save",routes.SaveRoutes)
	http.HandleFunc("/simul_search",simulSearchTpl)
	http.HandleFunc("/do_simul_search",routes.DoSimulSearch)
	http.HandleFunc("/",home)

	http.ListenAndServe(":80",nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := false
	isLoggedIn = auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	home_tpl.ExecuteTemplate(w, "home.html",data)
}
func askConfirm(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := false
	isLoggedIn = auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	auth_tpl.ExecuteTemplate(w, "ask_confirm_email.html",data)
}
func registerForm(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := false
	isLoggedIn = auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	auth_tpl.ExecuteTemplate(w, "register.html",data)
}
func loginForm(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := false
	isLoggedIn = auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	auth_tpl.ExecuteTemplate(w, "login.html",data)
}
func index(w http.ResponseWriter, req *http.Request){
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil{
		panic("Can't load env file")
	}
	//envファイルからAPI key取得
	apiKey := os.Getenv("MAP_API_KEY")

	isLoggedIn := false
	isLoggedIn = auth.IsLoggedIn(req)
	data := map[string]interface{}{"apiKey":apiKey,"isLoggedIn":isLoggedIn}
	tpl.ExecuteTemplate(w, "place_and_direction_improve.html", data)
}

func simulSearchTpl(w http.ResponseWriter, req *http.Request) {
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil{
		panic("Can't load env file")
	}
	//envファイルからAPI key取得
	apiKey := os.Getenv("MAP_API_KEY")
	isLoggedIn := false
	isLoggedIn = auth.IsLoggedIn(req)
	nineIterator := []int {1,2,3,4,5,6,7,8,9}
	data := map[string]interface{}{"apiKey":apiKey,"isLoggedIn":isLoggedIn,"nineIterator":nineIterator}
	simul_search_tpl.ExecuteTemplate(w, "simul_search.html",data)
}