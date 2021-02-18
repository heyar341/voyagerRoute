package routes

import (
	//envファイル操作用のパッケージ
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
	"html"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func DoSimulSearch(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Receive request")
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil {
		log.Println("Can't load env file")
	}
	//envファイルからAPI key取得
	apiKey := os.Getenv("MAP_API_KEY")

	//API呼び出しクライアントを作成
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(10))
	check(err)

	//目的地以外のパラメータは共通なので、mapに値を入れて、simulSerachメソッドに引数として渡す。
	reqParam := map[string]string{
		"origin":        "",
		"departureTime": "",
		"arrivalTime":   "",
		"mode":          "",
		"avoid":         "",
	}

	//出発地
	origin := html.EscapeString(req.FormValue("origin"))
	reqParam["origin"] = origin
	//出発時間はデフォルトで現在時刻に設定
	defDepartureTime := strconv.Itoa(int(time.Now().Unix()))
	reqParam["departureTime"] = defDepartureTime

	//時間指定する場合
	departureTime := html.EscapeString(req.FormValue("departureTime"))
	arrivalTime := html.EscapeString(req.FormValue("arrivalTime"))
	if departureTime != "" {
		reqParam["departureTime"] = departureTime
	} else if arrivalTime != "" {
		reqParam["departureTime"] = ""
		reqParam["arrivalTime"] = arrivalTime
	}

	//徒歩、公共交通機関、自動車を選択
	mode := html.EscapeString(req.FormValue("mode"))
	reqParam["mode"] = mode
	//自動車選択時、有料道路などを含めないオプション
	avoid := html.EscapeString(req.FormValue("avoid"))
	reqParam["avoid"] = avoid

	//同時検索
	for i := 1; i < 10; i++ {
		destination := html.EscapeString(req.FormValue("destination" + strconv.Itoa(i)))
		if destination == "" {
			break
		}
		disntance, duration := simulSearch(client, destination, reqParam)
		fmt.Print(disntance, duration)
		//同期か非同期か決まってから処理を決定
	}
	return
}

func simulSearch(client *maps.Client, destination string, reqParam map[string]string) (string, int) {
	//requestの変数宣言
	SearchReq := &maps.DirectionsRequest{
		Language:    "ja",
		Region:      "JP",
		Origin:      reqParam["origin"],
		Destination: destination,
		//過去のデータから予想される最適な所要時間を返すよう設定
		TrafficModel: maps.TrafficModelBestGuess,
	}

	//オプション指定されている場合、SearchReqにその値を入れる
	if reqParam["departureTime"] != "" {
		SearchReq.DepartureTime = reqParam["departureTime"]
	} else if reqParam["arrivalTime"] != "" {
		SearchReq.DepartureTime = ""
		SearchReq.ArrivalTime = reqParam["arrivalTime"]
	}

	if reqParam["mode"] != "" {
		lookupMode(reqParam["mode"], SearchReq)
	}
	if reqParam["avoid"] != "" {
		lookupAvoid(reqParam["avoid"], SearchReq)
	}

	//ルートを取得
	routes, _, err := client.Directions(context.Background(), SearchReq)
	check(err)

	//pretty.Println(routes[0].Legs[0].Distance)
	//pretty.Println(int(routes[0].Legs[0].Duration.Minutes()))
	return routes[0].Legs[0].Distance.HumanReadable, int(routes[0].Legs[0].Duration.Minutes())
}

//エラーチェック
func check(err error) {
	if err != nil {
		log.Printf("fatal error: %s", err)
	}
}

//ここより下、オプション設定メソッド
//徒歩、運転、乗り換えのモード選択
func lookupMode(mode string, r *maps.DirectionsRequest) {
	switch mode {
	case "driving":
		r.Mode = maps.TravelModeDriving
	case "walking":
		r.Mode = maps.TravelModeWalking
	case "bicycling":
		r.Mode = maps.TravelModeBicycling
	case "transit":
		r.Mode = maps.TravelModeTransit
	case "":
		// ignore
	default:
		log.Fatalf("Unknown mode '%s'", mode)
	}
}

//trueの場合複数ルートを探し、falseの場合１つのルートのみ返す
func lookupAlternatives(alternatives string, r *maps.DirectionsRequest) {
	if alternatives == "true" {
		r.Alternatives = true
	} else {
		r.Alternatives = false
	}
}

//乗り換え数の少なさを優先するか、歩行距離の短さを優先するか選択
func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DirectionsRequest) {
	switch transitRoutingPreference {
	case "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	case "":
		// ignore
	default:
		log.Fatalf("Unknown transit routing preference %s", transitRoutingPreference)
	}
}

//最速時間、過去のデータからの最適予測時間、最も遅い場合の予測のどれか選択
func lookupTrafficModel(trafficModel string, r *maps.DirectionsRequest) {
	switch trafficModel {
	case "optimistic":
		r.TrafficModel = maps.TrafficModelOptimistic
	case "best_guess":
		r.TrafficModel = maps.TrafficModelBestGuess
	case "pessimistic":
		r.TrafficModel = maps.TrafficModelPessimistic
	case "":
		// ignore
	default:
		log.Fatalf("Unknown traffic mode %s", trafficModel)
	}
}

//有料道路、高速道路、フェリーを除外する場合選択
func lookupAvoid(avoid string, r *maps.DirectionsRequest) {
	for _, a := range strings.Split(avoid, "|") {
		switch a {
		case "tolls":
			r.Avoid = append(r.Avoid, maps.AvoidTolls)
		case "highways":
			r.Avoid = append(r.Avoid, maps.AvoidHighways)
		case "ferries":
			r.Avoid = append(r.Avoid, maps.AvoidFerries)
		default:
			log.Fatalf("Unknown avoid restriction %s", a)
		}
	}
}

//交通手段を選択
func lookupTransitMode(transitMode string, r *maps.DirectionsRequest) {
	for _, t := range strings.Split(transitMode, "|") {
		switch t {
		case "bus":
			r.TransitMode = append(r.TransitMode, maps.TransitModeBus)
		case "subway":
			r.TransitMode = append(r.TransitMode, maps.TransitModeSubway)
		case "train":
			r.TransitMode = append(r.TransitMode, maps.TransitModeTrain)
		case "tram":
			r.TransitMode = append(r.TransitMode, maps.TransitModeTram)
		case "rail":
			r.TransitMode = append(r.TransitMode, maps.TransitModeRail)
		}
	}
}
