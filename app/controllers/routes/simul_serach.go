package routes

import (
	"app/controllers/envhandler"
	"app/controllers/routes/simuloptions"
	"context"
	"encoding/json"
	"fmt"
	"googlemaps.github.io/maps"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Resp struct {
	Field map[string][]string `json:"resp"`
}

type SimulSearchRequest struct {
	Origin        string            `json:"origin"`
	Destinations  map[string]string `json:"destinations"`
	Mode          string            `json:"mode"`
	DepartureTime string            `json:"departure_time"`
	Avoid         string            `json:"avoid"`
}

func DoSimulSearch(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Receive request")
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")

	//API呼び出しクライアントを作成
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(10))
	check(err)

	//目的地以外のパラメータは共通なので、mapに値を入れて、simulSerachメソッドに引数として渡す。
	reqParam := map[string]string{
		"origin":        "",
		"departureTime": "",
		"mode":          "",
		"avoid":         "",
	}

	//requestのフィールドを保存する変数
	var reqFields SimulSearchRequest
	body, _ := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &reqFields)
	if err != nil {
		http.Error(w, "aa", http.StatusInternalServerError)
	}

	//出発地
	origin := "place_id:" + reqFields.Origin
	reqParam["origin"] = origin

	//徒歩、公共交通機関、自動車を選択
	mode := reqFields.Mode
	reqParam["mode"] = mode
	//自動車選択時、有料道路などを含めないオプション
	if mode == "transit" {
		//時間指定する場合
		departureTime := reqFields.DepartureTime
		reqParam["departureTime"] = departureTime
	} else if mode == "driving" {
		//運転時オプションを指定する場合
		avoid := reqFields.Avoid
		reqParam["avoid"] = avoid
	}

	//検索結果を入れるmap
	simulRoutes := map[string][]string{}
	//同時検索
	for i := 1; i < 10; i++ {
		destination := reqFields.Destinations[strconv.Itoa(i)]
		if destination == "" {
			continue
		}
		destination = "place_id:" + destination
		disntance, duration := simulSearch(client, destination, reqParam)
		//エラーもしくは検索結果がない場合
		if disntance == "" && duration == 0 {
			simulRoutes[strconv.Itoa(i)] = []string{"検索結果なし", "検索結果なし"}
		} else {
			var durationResp string
			//１時間以上の場合、〜時間〜分にフォーマットを変える
			if duration >= 60 {
				durationResp = strconv.Itoa(int(duration/60)) + "時間" + strconv.Itoa(int(duration%60)) + "分"
			} else {
				durationResp = strconv.Itoa(duration) + "分"
			}
			simulRoutes[strconv.Itoa(i)] = []string{disntance, durationResp}
		}
	}

	//レスポンスを作成
	resp := Resp{
		Field: simulRoutes,
	}
	w.Header().Set("Content-Type", "application/json")
	respJson, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "問題が発生しました。もう一度操作しなおしてください", http.StatusInternalServerError)
	}
	w.Write(respJson)
}

func simulSearch(client *maps.Client, destination string, reqParam map[string]string) (string, int) {
	//requestの変数宣言
	SearchReq := &maps.DirectionsRequest{
		Language:    "ja",
		Region:      "JP",
		Origin:      reqParam["origin"],
		Destination: destination,
		//出発時間はデフォルトで現在時刻に設定
		DepartureTime: strconv.Itoa(int(time.Now().Unix())),
		//過去のデータから予想される最適な所要時間を返すよう設定
		TrafficModel: maps.TrafficModelBestGuess,
	}

	//オプション指定されている場合、SearchReqにその値を入れる
	if reqParam["departureTime"] != "" {
		t, _ := time.Parse(time.RFC3339, reqParam["departureTime"])
		SearchReq.DepartureTime = strconv.Itoa(int(t.Unix()))
	}
	if reqParam["mode"] != "" {
		simuloptions.LookupMode(reqParam["mode"], SearchReq)
	}
	if reqParam["avoid"] != "" {
		simuloptions.LookupAvoid(reqParam["avoid"], SearchReq)
	}

	//ルートを取得
	routes, _, err := client.Directions(context.Background(), SearchReq)
	check(err)
	if err != nil {
		return "", 0
	}
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
