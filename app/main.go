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
	//"fmt"
	//"log"
	"net/http"
	//"os"
	"app/direction"

)

func main() {

	http.Handle("/favicon.ico", http.NotFoundHandler())
	//Direction API
	http.HandleFunc("/route_search",searchRoute)
	http.ListenAndServe(":80",nil)
}

func searchRoute(w http.ResponseWriter, req *http.Request)  {
	routes := direction.SearchRoute(req)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}