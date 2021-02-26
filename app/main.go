// Copyright 2015 Google Inc. All Rights Reserved.
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"app/controllers/auth"
	"app/controllers/middleware"
	"app/controllers/mypage"
	"app/controllers/routes"
	"fmt"
	"html/template"
	"net/http"
)

var home_tpl *template.Template

func init() {
	home_tpl = template.Must(template.Must(template.ParseGlob("templates/home/home.html")).ParseGlob("templates/includes/*.html"))
}

func main() {

	fmt.Println("App started")
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	//Authentication
	http.HandleFunc("/register_form/", middleware.Auth(auth.RegisterForm))
	http.HandleFunc("/check_email", auth.EmailIsAvailable)
	http.HandleFunc("/register", middleware.RegisterValidator(auth.Register))
	http.HandleFunc("/login_form/", middleware.Auth(auth.LoginForm))
	http.HandleFunc("/login", middleware.LoginValidator(auth.Login))
	http.HandleFunc("/confirm_register/", auth.ConfirmRegister)
	http.HandleFunc("/ask_confirm/", middleware.Auth(auth.AskConfirmEmail))
	http.HandleFunc("/logout", auth.Logout)

	//Direction API
	http.HandleFunc("/multi_search", middleware.Auth(routes.MultiSearchTpl))
	http.HandleFunc("/routes_save", middleware.Auth(middleware.SaveRoutesValidator(routes.SaveRoutes)))
	http.HandleFunc("/simul_search", middleware.Auth(routes.SimulSearchTpl))
	http.HandleFunc("/do_simul_search", routes.DoSimulSearch)
	http.HandleFunc("/show_route/", middleware.Auth(routes.ShowAndEditRoutesTpl))
	http.HandleFunc("/update_route", routes.UpdateRoute)

	http.HandleFunc("/", middleware.Auth(home))

	http.HandleFunc("/mypage/show_routes", middleware.Auth(mypage.ShowAllRoutes))
	http.HandleFunc("/mypage", middleware.Auth(mypage.ShowMypage))

	http.ListenAndServe(":80", nil)
}

func home(w http.ResponseWriter, req *http.Request) {
	data := req.Context().Value("data").(map[string]interface{})
	data["msg"] = req.URL.Query().Get("msg")
	data["success"] = req.URL.Query().Get("success")
	home_tpl.ExecuteTemplate(w, "home.html", data)
}