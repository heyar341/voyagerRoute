package reqvalidator

import (
	"app/controllers/simulsearch"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func SimulSearchValidator(DoSimulSearch http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//リクエストメソッドについて確認
		if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
			http.Error(w, "リクエスト方法が不正です。", http.StatusBadRequest)
			log.Printf("Someone sended data not from simul_search page")
			return
		}
		//requestのフィールドを保存する変数
		var reqParams simulsearch.SimulParams
		body, _ := ioutil.ReadAll(req.Body)
		err := json.Unmarshal(body, &reqParams)
		if err != nil {
			http.Error(w, "リクエストデータに不備があります。", http.StatusBadRequest)
			log.Printf("Error while json marshaling simulSearch request: %v", err)
			return
		}

		//出発地のバリデーション
		if reqParams.Origin == "" {
			http.Error(w, "出発地を入力してください。", http.StatusBadRequest)
			return
		}
		//place_id:を追加
		reqParams.Origin = "place_id:" + reqParams.Origin

		//目的地のバリデーション
		for i := 1; i < 10; i++ {
			if reqParams.Destinations[strconv.Itoa(i)] == "" {
				continue
			}
			reqParams.Destinations[strconv.Itoa(i)] = "place_id:" + reqParams.Destinations[strconv.Itoa(i)]
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqParams", reqParams)
		DoSimulSearch.ServeHTTP(w, req.WithContext(ctx))
	}
}
