package simulsearch

import (
	"app/internal/customerr"
	"app/internal/envhandler"
	"app/internal/errormsg"
	"app/model"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"googlemaps.github.io/maps"
)

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
		log.Printf("Unknown mode '%s'", mode)
	}
}

//lookupAlternatives sets Alternatives field value in DirectionsRequest.
func lookupAlternatives(alternatives string, r *maps.DirectionsRequest) {
	if alternatives == "true" {
		r.Alternatives = true
	} else {
		r.Alternatives = false
	}
}

//lookupTransitRoutingPreference sets TransitRoutingPreference field value in DirectionsRequest.
func lookupTransitRoutingPreference(transitRoutingPreference string, r *maps.DirectionsRequest) {
	switch transitRoutingPreference {
	case "fewer_transfers":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceFewerTransfers
	case "less_walking":
		r.TransitRoutingPreference = maps.TransitRoutingPreferenceLessWalking
	case "":
		// ignore
	default:
		log.Printf("Unknown transit routing preference %s", transitRoutingPreference)
	}
}

//lookupTrafficModel sets TrafficModel field value in DirectionsRequest.
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
		log.Printf("Unknown traffic mode %s", trafficModel)
	}
}

//lookupAvoid sets Avoid field value in DirectionsRequest.
func lookupAvoid(avoid []string, r *maps.DirectionsRequest) {
	if len(avoid) == 0 {
		return
	}
	for _, a := range avoid {
		switch a {
		case "tolls":
			r.Avoid = append(r.Avoid, maps.AvoidTolls)
		case "highways":
			r.Avoid = append(r.Avoid, maps.AvoidHighways)
		case "ferries":
			r.Avoid = append(r.Avoid, maps.AvoidFerries)
		default:
			log.Printf("Unknown avoid restriction %s", a)
		}
	}
}

//lookupTransitMode sets TransitMode field value in DirectionsRequest.
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

type timeZoneResp struct {
	SummerTimeOffset int    `json:"dstOffset"` //サマータイム時のオフセット
	RawOffset        int    `json:"rawOffset"` //通常時のオフセット
	Status           string `json:"status"`
	TimeZoneID       string `json:"timeZoneId"`
	TimeZoneName     string `json:"timeZoneName"`
}

var timeZoneAPIURL string = "https://maps.googleapis.com/maps/api/timezone/json?location="

//getTimeZoneOffset fetches timezone offset from Google Maps API TimeZone API in second unit.
//the returned value from TimeZone API is converted to hour unit and formatted as RFC3339.
func getTimeZoneOffset(lat, lng string) (string, error) {
	apiKey, err := envhandler.GetEnvVal("TIMEZONE_API_KEY")
	if err != nil {
		return "", err
	}
	//timezone API用URL
	reqURL := timeZoneAPIURL + lat + "," + lng + "&timestamp=" +
		strconv.Itoa(int(time.Now().Unix())) + "&key=" + apiKey

	resp, err := http.Get(reqURL)
	if err != nil {
		log.Printf("Error while getting timezone response: %v", err)
		return "", err
	}
	//responseのフィールドを保存する変数
	var tZResp timeZoneResp
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tZResp)
	if err != nil {
		log.Printf("Error while json unmarshaling timezone response: %v", err)
		return "", err
	}
	err = resp.Body.Close()

	offset := tZResp.RawOffset       //seconds
	offsetHour := int(offset / 3600) //hours

	var strOffset string
	switch {
	case offsetHour == 0:
		strOffset = "Z" //YTC
	case offsetHour > 0:
		//UTCより進んでいる場所
		strOffset = "0" + strconv.Itoa(offsetHour)
		strOffset = "+" + strOffset[len(strOffset)-2:] + ":00" //1桁なら+0(number):00, 2桁なら+1or2(number):00
	case offsetHour < 0:
		//UTCより遅れている場所
		strOffset = "0" + strconv.Itoa(-offsetHour)
		strOffset = "-" + strOffset[len(strOffset)-2:] + ":00"
	}
	return strOffset, nil

}

//getDataFromAPI fetches distance and duration from Google Maps API Directions API.
// Request is initiated with default values. According to specified options, some
//fields(Mode, DepartureTime TrafficModel) are set.
//distance is human readable format like ~m or ~km, and distance is integer format in minute unit.
func getDataFromAPI(client *maps.Client, destination string, reqParam *SimulParams) (string, int) {
	//requestの変数宣言
	searchReq := &maps.DirectionsRequest{
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

	lookupMode(reqParam.Mode, searchReq)
	switch reqParam.Mode {
	case "transit":
		//オプション指定されている場合、searchReqにその値を入れる
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
			searchReq.DepartureTime = strconv.Itoa(int(t.Unix()))
		}
	case "driving":
		lookupAvoid(reqParam.Avoid, searchReq)
	}
	//ルートを取得
	routes, _, err := client.Directions(context.Background(), searchReq)
	if err != nil || len(routes) == 0 {
		return "", 0
	}
	distance := routes[0].Legs[0].Distance.HumanReadable
	duration := int(routes[0].Legs[0].Duration.Minutes())
	return distance, duration
}

//convertDurationToStr converts duration in minute unit to datetime format in string.
func convertDurationToStr(duration int) string {
	var d string
	//１時間以上の場合、〜時間〜分にフォーマットを変える
	switch {
	case duration >= 24*60:
		d = strconv.Itoa(int(duration/(24*60))) + "日" +
			strconv.Itoa(int(duration%(24*60)/60)) + "時間" +
			strconv.Itoa(int(duration%(24*60)%60)) + "分"
	case duration >= 60:
		d = strconv.Itoa(int(duration/60)) + "時間" +
			strconv.Itoa(int(duration%60)) + "分"
	default:
		d = strconv.Itoa(duration) + "分"
	}
	return d
}

type searchRoute struct {
	reqParams    SimulParams
	apiKey       string
	client       *maps.Client
	destinations map[string]model.DestinationData
	err          error
}

//getReqParamFromCtx fetch request parameters from context.
func (s *searchRoute) getReqParamFromCtx(req *http.Request) {
	//Validation後の炉クエストパラメータを取得
	reqParams, ok := req.Context().Value("reqParams").(SimulParams)
	if !ok {
		s.err = customerr.BaseErr{
			Op:  "get reqParams from context",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while getting reqParams from context: %v", ok),
		}
		return
	}
	s.reqParams = reqParams
}

//getAPIkeyFromEnv fetch Google Maps API key from .env file.
func (s *searchRoute) getAPIkeyFromEnv() {
	//envファイルからAPIキー取得
	apiKey, err := envhandler.GetEnvVal("MAP_API_KEY")
	if err != nil {
		s.err = customerr.BaseErr{
			Op:  "get APIkey from env file",
			Msg: errormsg.SomethingBad,
			Err: fmt.Errorf("error while getting APIkey from env file: %v", err),
		}
		return
	}
	s.apiKey = apiKey
}

//genNewClient generate new Google Maps API client.
func (s *searchRoute) genNewClient() {
	//API呼び出しクライアントを作成
	client, err := maps.NewClient(maps.WithAPIKey(s.apiKey), maps.WithRateLimit(10))
	if err != nil {
		s.err = customerr.BaseErr{
			Op:  "construct map API client",
			Msg: errormsg.TriAgain,
			Err: fmt.Errorf("couldn't use Directions API: %v", err),
		}
		return
	}
	s.client = client

}

//executeSearch sets place_id, distance and duration value of route from 1 to 9.
//place id is prefixed with "place_id:"(9 characters) and it's used in frontend
//without prefix, so it's formatted to without prefix value.
//if route isn't found, values will be "検索結果なし".
func (s *searchRoute) executeSearch() {
	//同時検索
	wg := sync.WaitGroup{}
	mux := sync.Mutex{}
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			mux.Lock()
			destination := s.reqParams.Destinations[strconv.Itoa(i)]
			mux.Unlock()
			if destination == "" {
				return
			}
			distance, duration := getDataFromAPI(s.client, destination, &s.reqParams)
			//エラーもしくは検索結果がない場合
			if distance == "" && duration == 0 {
				mux.Lock()
				s.destinations[strconv.Itoa(i)] = model.DestinationData{
					PlaceId:  destination[9:],
					Distance: "検索結果なし",
					Duration: "検索結果なし",
				}
				mux.Unlock()
			} else {
				d := convertDurationToStr(duration)
				mux.Lock()
				s.destinations[strconv.Itoa(i)] = model.DestinationData{
					PlaceId:  destination[9:],
					Distance: distance,
					Duration: d,
				}
				mux.Unlock()
			}
		}(i)
	}
	wg.Wait()
}

func Search(w http.ResponseWriter, req *http.Request) {
	var s searchRoute
	s.getReqParamFromCtx(req)
	s.getAPIkeyFromEnv()
	s.genNewClient()
	if s.err != nil {
		e := s.err.(customerr.BaseErr)
		http.Error(w, e.Msg, http.StatusInternalServerError)
		log.Println(e.Err)
		return
	}
	s.destinations = make(map[string]model.DestinationData)
	s.executeSearch()

	//レスポンスを作成
	w.Header().Set("Content-Type", "application/json")
	respJson, err := json.Marshal(s.destinations)
	if err != nil {
		http.Error(w, "問題が発生しました。もう一度操作しなおしてください", http.StatusInternalServerError)
		return
	}
	w.Write(respJson)
}
