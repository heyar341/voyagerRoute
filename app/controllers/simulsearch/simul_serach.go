package simulsearch

import (
	"app/controllers/envhandler"
	"context"
	"encoding/json"
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

type Params struct {
	Origin        string            `json:"origin"`
	Destinations  map[string]string `json:"destinations"`
	Mode          string            `json:"mode"`
	DepartureTime string            `json:"departure_time"`
	Avoid         string            `json:"avoid"`
}

func DoSimulSearch(w http.ResponseWriter, req *http.Request) {
	//envファイルからAPIキー取得
	apiKey := envhandler.GetEnvVal("MAP_API_KEY")

	//API呼び出しクライアントを作成
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(10))
	if err != nil {
		http.Error(w, "APIが使用できません。しばらく経ってからもう一度お試しください。", http.StatusInternalServerError)
		log.Printf("Couldn't use Directions API :%v", err)
		return
	}

	//requestのフィールドを保存する変数
	var reqParams Params
	body, _ := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &reqParams)
	if err != nil {
		http.Error(w, "リクエストデータに不備があります。", http.StatusBadRequest)
		log.Printf("Error while json marshaling simulSearch request: %v", err)
	}

	reqParams.Origin = "place_id:" + reqParams.Origin
	//検索結果を入れるmap
	simulRoutes := map[string][]string{}
	//同時検索
	for i := 1; i < 10; i++ {
		destination := reqParams.Destinations[strconv.Itoa(i)]
		if destination == "" {
			continue
		}
		destination = "place_id:" + destination
		disntance, duration := simulSearch(client, destination, &reqParams)
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

//google maps Directions APIを使用して、距離と所要時間お取得する関数
func simulSearch(client *maps.Client, destination string, reqParam *Params) (string, int) {
	//requestの変数宣言
	SearchReq := &maps.DirectionsRequest{
		Language:    "ja",
		Region:      "JP",
		Origin:      reqParam.Origin,
		Destination: destination,
		//出発時間はデフォルトで現在時刻に設定
		DepartureTime: strconv.Itoa(int(time.Now().Unix())),
		//過去のデータから予想される最適な所要時間を返すよう設定
		TrafficModel: maps.TrafficModelBestGuess,
	}

	//オプション指定されている場合、SearchReqにその値を入れる
	if reqParam.DepartureTime != "" {
		t, _ := time.Parse(time.RFC3339, reqParam.DepartureTime)
		SearchReq.DepartureTime = strconv.Itoa(int(t.Unix()))
	}
	if reqParam.Mode != "" {
		lookupMode(reqParam.Mode, SearchReq)
	}
	if reqParam.Avoid != "" {
		lookupAvoid(reqParam.Avoid, SearchReq)
	}

	//ルートを取得
	routes, _, err := client.Directions(context.Background(), SearchReq)
	if err != nil {
		return "", 0
	}
	return routes[0].Legs[0].Distance.HumanReadable, int(routes[0].Legs[0].Duration.Minutes())
}
