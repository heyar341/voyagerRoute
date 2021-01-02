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
)
var tpl *template.Template

type ResponseSample struct {
	Field1 string
	Field2 []string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/route_search/*"))
}

func main() {

	http.Handle("/favicon.ico", http.NotFoundHandler())
	//Direction API
	http.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("./templates"))))
	http.HandleFunc("/show_map",index)
	http.HandleFunc("/routes_save",saveRoutes)

	http.ListenAndServe(":80",nil)
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
}