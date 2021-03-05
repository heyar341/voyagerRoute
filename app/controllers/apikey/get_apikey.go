package apikey

import (
	"app/envhandler"
	"io/ioutil"
	"log"
	"net/http"
)

func GetApiKey(w http.ResponseWriter, req *http.Request) {
	apiKey, err := envhandler.GetEnvVal("MAP_API_KEY")
	if err != nil || apiKey == "" {
		http.Error(w, "エラーが発生しました。", http.StatusBadRequest)
		log.Println("Couldn't get API Key from env file")
		return
	}

	//Google Maps APIのJavascriptファイルを取得するURLを生成
	baseURL := "https://maps.googleapis.com/maps/api/js?key="
	optsURL := "&libraries=places&v=weekly&language=ja"
	reqURL := baseURL + apiKey + optsURL

	//Javascriptファイルを取得
	apiRes, err := http.Get(reqURL)
	if err != nil {
		http.Error(w, "エラーが発生しました。", http.StatusBadRequest)
		log.Printf("Couldn't get file response from Google Maps API seervice %v", err)
		return
	}
	//HTTPレスポンスとして渡せるよう[]byte型にする
	srcJS, err := ioutil.ReadAll(apiRes.Body)
	if err != nil {
		http.Error(w, "エラーが発生しました。", http.StatusBadRequest)
		log.Printf("Error while reading file returned from Google Maps API service: %v", err)
		return
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "text/javascript;charset=UTF-8")
	_, err = w.Write(srcJS)
	if err != nil {
		log.Printf("Error while returning javascript file for Google Maps API service: %v", err)
		return
	}
}
