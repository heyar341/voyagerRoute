package multiroute

import (
	"app/controllers/envhandler"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type TimeZoneReq struct {
	Lat      string `json:"lat"`       //緯度
	Lng      string `json:"lng"`       //経度
	UnixTime string `json:"unix_time"` //現在時刻のUnix表記
}

type TimeZoneResp struct {
	SummerTimeOffset int    `json:"dstOffset"` //サマータイム時のオフセット
	RawOffset        int    `json:"rawOffset"` //通常時のオフセット
	Status           string `json:"status"`
	//TimeZoneID       string `json:"timeZoneId"`
	//TimeZoneName     string `json:"timeZoneName"`
}

//TimeZone API公式ドキュメント：https://developers.google.com/maps/documentation/timezone/get-started

func GetTimezone(w http.ResponseWriter, req *http.Request) {
	//リクエストメソッドについて確認
	if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
		http.Error(w, "リクエストメソッドが不正です。", http.StatusBadRequest)
		log.Printf("Someone sended data not from multi_search page for TimeZone data")
		return
	}

	apiKey, err := envhandler.GetEnvVal("TIMEZONE_API_KEY")
	if err != nil {
		msg := "エラーが発生しました。現在サービスをご利用いただけません。"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	//requestのフィールドを保存する変数
	var tZReq TimeZoneReq
	body, _ := ioutil.ReadAll(req.Body)
	err = json.Unmarshal(body, &tZReq)
	if err != nil {
		http.Error(w, "リクエストデータに不備があります。", http.StatusBadRequest)
		log.Printf("Error while json unmarshaling timezone request: %v", err)
		return
	}
	err = req.Body.Close()

	//timezone API用URL
	reqURL := "https://maps.googleapis.com/maps/api/timezone/json?location=" +
		tZReq.Lat + "," + tZReq.Lng + "&timestamp=" + tZReq.UnixTime + "&key=" + apiKey

	resp, err := http.Get(reqURL)
	if err != nil {
		http.Error(w, "データの取得に失敗しました。", http.StatusInternalServerError)
		log.Printf("Error while json unmarshaling timezone response: %v", err)
		return
	}
	//requestのフィールドを保存する変数
	var tZResp TimeZoneResp
	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tZResp)
	if err != nil {
		http.Error(w, "データの取得に失敗しました。", http.StatusInternalServerError)
		log.Printf("Error while json unmarshaling timezone response: %v", err)
		return
	}
	err = resp.Body.Close()

	type OffsetJSON struct {
		RawOffset string `json:"rawOffset"`
	}
	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	offsetInfo := OffsetJSON{RawOffset: strconv.Itoa(tZResp.RawOffset)}
	respJson, err := json.Marshal(offsetInfo)
	if err != nil {
		http.Error(w, "データの取得に失敗しました。", http.StatusInternalServerError)
		log.Printf("Error while json marshaling timezone data: %v", err)
		return
	}
	w.Write(respJson)
}
