// Copyright 2015 Google Inc. All Rights Reserved.
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"app/controllers/auth"
	"app/controllers/envhandler"
	"app/controllers/middleware"
	"app/controllers/mypages"
	"app/controllers/routes"
	"fmt"
	"html/template"
	"net/http"
)

var tpl, auth_tpl, home_tpl, simul_search_tpl, mypage_tpl, show_route_tpl *template.Template

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
	http.HandleFunc("/register_form/", middleware.Auth(registerForm))
	http.HandleFunc("/check_email", auth.EmailIsAvailable)
	http.HandleFunc("/register", middleware.RegisterValidator(auth.Register))
	http.HandleFunc("/login_form/", middleware.Auth(loginForm))
	http.HandleFunc("/login", middleware.LoginValidator(auth.Login))
	http.HandleFunc("/confirm_register/", auth.ConfirmRegister)
	http.HandleFunc("/ask_confirm/", middleware.Auth(askConfirm))
	http.HandleFunc("/logout", auth.Logout)

	//Direction API
	http.HandleFunc("/multi_search", middleware.Auth(index))
	http.HandleFunc("/routes_save", middleware.Auth(middleware.SaveRoutesValidator(routes.SaveRoutes)))
	http.HandleFunc("/simul_search", middleware.Auth(simulSearchTpl))
	http.HandleFunc("/do_simul_search", routes.DoSimulSearch)
	http.HandleFunc("/show_route/", middleware.Auth(showAndEditRoutes))
	http.HandleFunc("/update_route", routes.UpdateRoute)

	http.HandleFunc("/", middleware.Auth(home))

	http.HandleFunc("/mypage/show_routes", middleware.Auth(showRoutes))
	http.HandleFunc("/mypage", middleware.Auth(mypage))

	http.ListenAndServe(":80", nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	data["msg"] = req.URL.Query().Get("msg")
	data["success"] = req.URL.Query().Get("success")
	home_tpl.ExecuteTemplate(w, "home.html", data)
}

func mypage(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	user := req.Context().Value("user").(middleware.UserData)
	data["userName"] = user.UserName
	mypage_tpl.ExecuteTemplate(w, "mypage.html", data)
}

func showRoutes(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	user := req.Context().Value("user").(middleware.UserData)
	titleNames := mypages.RouteTitles(user.ID)
	data["userName"] = user.UserName
	data["titles"] = titleNames
	mypage_tpl.ExecuteTemplate(w, "show_routes.html", data)
}
func askConfirm(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	auth_tpl.ExecuteTemplate(w, "ask_confirm_email.html", data)
}
func registerForm(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	data["qParams"] = req.URL.Query()
	auth_tpl.ExecuteTemplate(w, "register.html", data)
}
func loginForm(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	data["qParams"] = req.URL.Query()
	auth_tpl.ExecuteTemplate(w, "login.html", data)
}
func index(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	data["apiKey"] = apiKey
	tpl.ExecuteTemplate(w, "multi_search.html", data)
}

func simulSearchTpl(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	nineIterator := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	data["apiKey"] = apiKey
	data["nineIterator"] = nineIterator
	simul_search_tpl.ExecuteTemplate(w, "simul_search.html", data)
}

func showAndEditRoutes(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")
	data := req.Context().Value("data").(map[string]interface{})
	routeTitle := req.URL.Query().Get("route_title")
	user := req.Context().Value("user").(middleware.UserData)
	routeInfo := routes.GetRoute(w, routeTitle, user.ID)
	data["apiKey"] = apiKey
	data["routeInfo"] = routeInfo
	show_route_tpl.ExecuteTemplate(w, "multi_route_show.html", data)
}
