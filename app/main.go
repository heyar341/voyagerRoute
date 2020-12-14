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
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"app/direction"
	//envファイル操作用のパッケージ
	"github.com/joho/godotenv"

	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
)

var apiKey = flag.String("key", "", "API Key for using Google Maps API.")

func usageAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	fmt.Println("Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}

func main() {
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil{
		panic("Can't load env file")
	}
	*apiKey = os.Getenv("MAP_API_KEY")
	flag.Parse()

	var client *maps.Client
	var err error
	if *apiKey != "" {
		client, err = maps.NewClient(maps.WithAPIKey(*apiKey), maps.WithRateLimit(2))
	} else {
		usageAndExit("Please specify an API Key, or Client ID and Signature.")
	}
	check(err)
	//Direction API
	r := direction.MakeRequest(Origin,Destination,DepartureTime)
	routes, waypoints, err := client.Directions(context.Background(), r)
	check(err)
	pretty.Println(waypoints)
	pretty.Println(routes)
}



