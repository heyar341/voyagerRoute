// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains a simple command line tool for Directions API
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	//"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	//"app/direction"
	"html/template"
	"app/controllers/auth"

)
var route_tpl,auth_tpl *template.Template
func init() {
	route_tpl = template.Must(template.ParseGlob("templates/route_search/*"))
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
	route_tpl.ExecuteTemplate(w, "place_and_direction_improve.html", data)
}