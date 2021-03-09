package simulsearch

import (
	"app/envhandler"
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

//同時検索のリクエストパラメータ
type SimulParams struct {
	Origin        string            `json:"origin"`
	Destinations  map[string]string `json:"destinations"`
	Mode          string            `json:"mode"`
	DepartureTime string            `json:"departure_time"`
	LatLng        LatLng            `json:"latlng"`
	Avoid         []string          `json:"avoid"`
}

//SimulParamの緯度と経度
type LatLng struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type TimeZoneResp struct {
	SummerTimeOffset int    `json:"dstOffset"` //サマータイム時のオフセット
	RawOffset        int    `json:"rawOffset"` //通常時のオフセット
	Status           string `json:"status"`
	TimeZoneID       string `json:"timeZoneId"`
	TimeZoneName     string `json:"timeZoneName"`
}

func DoSimulSearch(w http.ResponseWriter, req *http.Request) {
	//Validation後の炉クエストパラメータを取得
	reqParams, ok := req.Context().Value("reqParams").(SimulParams)
	if !ok {
		http.Error(w, "リクエストパラメータに不備があります。", http.StatusInternalServerError)
		log.Printf("Error while getting reqParams from context: %v", ok)
		return
	}

	//envファイルからAPIキー取得
	apiKey, err := envhandler.GetEnvVal("MAP_API_KEY")
	if err != nil {
		http.Error(w, "エラーが発生しました。", http.StatusInternalServerError)
		return
	}
	//API呼び出しクライアントを作成
	client, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithRateLimit(10))
	if err != nil {
		http.Error(w, "APIが使用できません。しばらく経ってからもう一度お試しください。", http.StatusInternalServerError)
		log.Printf("Couldn't use Directions API :%v", err)
		return
	}

	//検索結果を入れるmap
	simulRoutes := map[string][]string{}

	//同時検索
	for i := 1; i < 10; i++ {
		destination := reqParams.Destinations[strconv.Itoa(i)]
		if destination == "" {
			continue
		}
		disntance, duration := simulSearch(client, destination, &reqParams)
		//エラーもしくは検索結果がない場合
		if disntance == "" && duration == 0 {
			simulRoutes[strconv.Itoa(i)] = []string{"検索結果なし", "検索結果なし"}
		} else {
			var durationResp string
			//１時間以上の場合、〜時間〜分にフォーマットを変える
			if duration >= 24*60 {
				durationResp = strconv.Itoa(int(duration/(24*60))) + "日" + strconv.Itoa(int(duration%(24*60)/60)) + "時間" +
					strconv.Itoa(int(duration%(24*60)%60)) + "分"
			} else if duration >= 60 {
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
func simulSearch(client *maps.Client, destination string, reqParam *SimulParams) (string, int) {
	//requestの変数宣言
	SearchReq := &maps.DirectionsRequest{
		Language:    "ja",
		Region:      "JP",
		Origin:      reqParam.Origin,
		Destination: destination,
		Mode:        maps.TravelModeWalking,
		//出発時間はデフォルトで現在時刻に設定
		DepartureTime: strconv.Itoa(int(time.Now().Unix())),
		//過去のデータから予想される最適な所要時間を返すよう設定
		TrafficModel: maps.TrafficModelBestGuess,
	}

	if reqParam.Mode == "trasit" {
		lookupMode(reqParam.Mode, SearchReq)
		//オプション指定されている場合、SearchReqにその値を入れる
		if reqParam.DepartureTime != "" {
			lat := reqParam.LatLng.Lat
			lng := reqParam.LatLng.Lng
			offset, err := getTimeZoneOffset(lat, lng)
			//オフセットを追加して、出発地のタイムゾーンの時間に合わせる
			if err != nil {
				return "", 0
			}
			specTime := reqParam.DepartureTime + offset
			t, err := time.Parse(time.RFC3339, specTime)
			if err != nil {
				log.Printf("Invalid specTime :%v", err)
				return "", 0
			}
			SearchReq.DepartureTime = strconv.Itoa(int(t.Unix()))
		}
	} else if reqParam.Mode == "driving" {
		lookupMode(reqParam.Mode, SearchReq)
		if len(reqParam.Avoid) > 0 {
			lookupAvoid(reqParam.Avoid, SearchReq)
		}
	} else {
		lookupMode(reqParam.Mode, SearchReq)
	}

	//ルートを取得
	routes, _, err := client.Directions(context.Background(), SearchReq)
	if err != nil {
		return "", 0
	}
	return routes[0].Legs[0].Distance.HumanReadable, int(routes[0].Legs[0].Duration.Minutes())
}

//緯度と経度からタイムゾーンオフセットを取得する関数
func getTimeZoneOffset(lat, lng string) (string, error) {
	apiKey, err := envhandler.GetEnvVal("TIMEZONE_API_KEY")
	if err != nil {
		return "", err
	}
	//timezone API用URL
	reqURL := "https://maps.googleapis.com/maps/api/timezone/json?location=" +
		lat + "," + lng + "&timestamp=" + strconv.Itoa(int(time.Now().Unix())) + "&key=" + apiKey

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Printf("Error while getting timezone response: %v", err)
		return "", err
	}
	//responseのフィールドを保存する変数
	var tZResp TimeZoneResp
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tZResp)
	if err != nil {
		log.Printf("Error while json unmarshaling timezone response: %v", err)
		return "", err
	}
	err = resp.Body.Close()
	offset := tZResp.RawOffset       //seconds
	offsetHour := int(offset / 3600) //hours
	if offsetHour == 0 {
		return "Z", nil //UTC
	} else if offsetHour > 0 {
		//UTCより進んでいる場所
		strOffset := "0" + strconv.Itoa(offsetHour)
		strOffset = "+" + strOffset[len(strOffset)-2:] + ":00" //1桁なら+0(number):00, 2桁なら+1or2(number):00
		return strOffset, nil
	} else {
		//UTCより遅れている場所
		strOffset := "0" + strconv.Itoa(-offsetHour)
		strOffset = "-" + strOffset[len(strOffset)-2:] + ":00"
		return strOffset, nil
	}
}
