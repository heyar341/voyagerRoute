// Copyright 2015 Google Inc. All Rights Reserved.
// Directions docs: https://developers.google.com/maps/documentation/directions/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	//"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	//"app/direction"
	"html/template"
	"app/controllers/auth"
)
var tpl *template.Template

type ResponseSample struct {
	Field1 string
	Field2 []string
}

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
	http.HandleFunc("/routes_save",saveRoutes)

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

func saveRoutes(w http.ResponseWriter, req *http.Request) {
	var Jsons interface{}
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body,&Jsons)
	if err != nil {
		http.Error(w, "aa", http.StatusInternalServerError)
	}

	resp := ResponseSample{"Heloo",[]string{"James","Bean","Bond"}}
	responeSample, err := json.Marshal(resp)

	fmt.Printf("%T",responeSample)
	if err != nil{
		http.Error(w,"問題が発生しました。もう一度操作しなおしてください",http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responeSample)
	route_tpl.ExecuteTemplate(w, "place_and_direction_improve.html", responeSample)
}