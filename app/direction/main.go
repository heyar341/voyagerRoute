package direction

import (
	"context"
	//envファイル操作用のパッケージ
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
	"os"
	"strings"
)

//requestの変数宣言
var Req = &maps.DirectionsRequest{
	Language:      "ja",
	Region:        "JP",
}

func SearchRoute(req *http.Request) []maps.Route{
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil{
		panic("Can't load env file")
	}
	//envファイルからAPI key取得
	apiKey := os.Getenv("MAP_API_KEY")
	var client *maps.Client
	var err error
	client, err = maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(2))
	check(err)
	//必須パラメータ
	Req.Origin = req.FormValue("origin")
	Req.Destination = req.FormValue("destination")
	Req.DepartureTime = req.FormValue("departureTime")
	Req.ArrivalTime = req.FormValue("arrivalTime")
	//オプションパラメータ
	if req.FormValue("mode") != "" { lookupMode(req.FormValue("mode"), Req) }
	if req.FormValue("waypoints") != "" { Req.Waypoints = strings.Split(req.FormValue("waypoints"), "|") }
	if req.FormValue("alternatives") != "" { lookupAlternatives(req.FormValue("alternatives"), Req) }
	if req.FormValue("transitRoutingPreference") != "" {lookupTransitRoutingPreference(req.FormValue("transitRoutingPreference"), Req)}
	if req.FormValue("trafficModel") != "" { lookupTrafficModel(req.FormValue("trafficModel"), Req) }
	if req.FormValue("avoid") != "" { lookupAvoid(req.FormValue("avoid"), Req) }
	if req.FormValue("transitMode") != "" { lookupTransitMode(req.FormValue("transitMode"), Req) }
	//ルート、中継地点を取得
	routes, waypoints, err := client.Directions(context.Background(), Req)
	check(err)
	pretty.Println(waypoints)
	pretty.Println(routes)
	return routes
}

//エラーチェック
func check(err error) {
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
}