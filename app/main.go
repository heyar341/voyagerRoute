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
	"encoding/json"
	"github.com/joho/godotenv"

	//"fmt"
	//"log"
	"net/http"
	"os"
	"app/direction"
	"html/template"

)
var tpl *template.Template
func init() {
	tpl = template.Must(template.ParseGlob("test/*"))
}

func main() {

	http.Handle("/favicon.ico", http.NotFoundHandler())
	//Direction API
	http.HandleFunc("/route_search",searchRoute)
	http.Handle("/test/", http.StripPrefix("/test", http.FileServer(http.Dir("./test"))))
	http.HandleFunc("/show_map",index)

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
	tpl.ExecuteTemplate(w, "place_and_direction.html", data)
}
func searchRoute(w http.ResponseWriter, req *http.Request)  {
	routes := direction.SearchRoute(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}