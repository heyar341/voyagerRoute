// Copyright 2015 Google Inc. All Rights Reserved.
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"fmt"
	"net/http"
	"html/template"
	"app/controllers/auth"
	"app/controllers/mypages"
	"app/controllers/routes"
	"app/controllers/envhandler"

)
var tpl, auth_tpl, home_tpl, simul_search_tpl, mypage_tpl,show_route_tpl *template.Template

func init() {
	tpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/search/*")).ParseGlob("templates/includes/*.html"))
	simul_search_tpl = template.Must(template.Must(template.ParseGlob("templates/simul_search/*")).ParseGlob("templates/includes/*.html"))
	auth_tpl = template.Must(template.Must(template.ParseGlob("templates/auth/*")).ParseGlob("templates/includes/*.html"))
	home_tpl = template.Must(template.Must(template.ParseGlob("templates/home/home.html")).ParseGlob("templates/includes/*.html"))
	mypage_tpl = template.Must(template.Must(template.ParseGlob("templates/mypage/*.html")).ParseGlob("templates/includes/*.html"))
	show_route_tpl = template.Must(template.Must(template.ParseGlob("templates/multi_search/show_and_edit/*")).ParseGlob("templates/includes/*.html"))

}

func main() {

	fmt.Println("App started")
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	//Authentication
	http.HandleFunc("/register_form/",registerForm)
	http.HandleFunc("/check_email",auth.EmailIsAvailable)
	http.HandleFunc("/register",auth.Register)
	http.HandleFunc("/login_form/",loginForm)
	http.HandleFunc("/login",auth.Login)
	http.HandleFunc("/confirm_register/",auth.ConfirmRegister)
	http.HandleFunc("/ask_confirm/",askConfirm)
	http.HandleFunc("/logout",auth.Logout)

	//Direction API
	http.HandleFunc("/multi_search",index)
	http.HandleFunc("/routes_save",routes.SaveRoutes)
	http.HandleFunc("/simul_search",simulSearchTpl)
	http.HandleFunc("/do_simul_search",routes.DoSimulSearch)
	http.HandleFunc("/show_route/",showAndEditRoutes)

	http.HandleFunc("/",home)

	http.HandleFunc("/mypage/show_routes", showRoutes)
	http.HandleFunc("/mypage", mypage)

	http.ListenAndServe(":80",nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	home_tpl.ExecuteTemplate(w, "home.html",data)
}

func mypage(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := auth.IsLoggedIn(req)
	userID, _ := auth.GetLoginUserID(req)
	userName, _ := auth.GetLoginUserName(userID)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn, "userName":userName}
	mypage_tpl.ExecuteTemplate(w, "mypage.html",data)
}

func showRoutes(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := auth.IsLoggedIn(req)
	userID, _ := auth.GetLoginUserID(req)
	userName, _ := auth.GetLoginUserName(userID)
	titleNames := mypages.RouteTitles(userID)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn, "userName":userName, "titles":titleNames}
	mypage_tpl.ExecuteTemplate(w, "show_routes.html",data)
}
func askConfirm(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	auth_tpl.ExecuteTemplate(w, "ask_confirm_email.html",data)
}
func registerForm(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	auth_tpl.ExecuteTemplate(w, "register.html",data)
}
func loginForm(w http.ResponseWriter, req *http.Request) {
	isLoggedIn := auth.IsLoggedIn(req)
	data := map[string]interface{}{"isLoggedIn":isLoggedIn}
	auth_tpl.ExecuteTemplate(w, "login.html",data)
}
func index(w http.ResponseWriter, req *http.Request){
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	isLoggedIn := auth.IsLoggedIn(req)
	data := map[string]interface{}{"apiKey":apiKey,"isLoggedIn":isLoggedIn}
	tpl.ExecuteTemplate(w, "multi_search.html", data)
}

func simulSearchTpl(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	isLoggedIn := auth.IsLoggedIn(req)
	nineIterator := []int {1,2,3,4,5,6,7,8,9}
	data := map[string]interface{}{"apiKey":apiKey,"isLoggedIn":isLoggedIn,"nineIterator":nineIterator}
	simul_search_tpl.ExecuteTemplate(w, "simul_search.html",data)
}

func showAndEditRoutes(w http.ResponseWriter, req *http.Request){
//envファイルからAPIキー取得
apiKey := envhandler.GetEnvVal("MAP_API_KEY")
isLoggedIn := auth.IsLoggedIn(req)
route_title := req.URL.Query().Get("route_title")
userID, _ := auth.GetLoginUserID(req)
routeInfo := routes.GetRoute(w,route_title, userID)
data := map[string]interface{}{"apiKey":apiKey,"isLoggedIn":isLoggedIn,"routeInfo":routeInfo}
show_route_tpl.ExecuteTemplate(w, "multi_route_show.html", data)
}