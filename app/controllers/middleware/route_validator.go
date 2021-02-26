package middleware

import (
	"app/controllers/routes"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func SaveRoutesValidator(SaveRoutes http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
			return
		}
		//requestのフィールドを保存する変数
		var reqFields routes.MultiSearchRequest
		body, _ := ioutil.ReadAll(req.Body)
		err := json.Unmarshal(body, &reqFields)
		if err != nil {
			http.Error(w, "入力に不正があります。", http.StatusBadRequest)
			log.Printf("Error while json marshaling: %v", err)
			return
		}

		if strings.ContainsAny(reqFields.Title, ".$") {
			http.Error(w, "ルート名にご使用いただけない文字が含まれています。", http.StatusBadRequest)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqFields", reqFields)
		SaveRoutes.ServeHTTP(w, req.WithContext(ctx))
	}
}

func UpdateRouteValidator(UpdateRoute http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
			return
		}
		//requestのフィールドを保存する変数
		type RouteUpdateRequest struct {
			ID            primitive.ObjectID     `json:"id" bson:"_id"`
			Title         string                 `json:"title" bson:"title"`
			PreviousTitle string                 `json:"previous_title" bson:"previous_title"`
			Routes        map[string]interface{} `json:"routes" bson:"routes"`
		}
		var reqFields RouteUpdateRequest
		body, _ := ioutil.ReadAll(req.Body)
		err := json.Unmarshal(body, &reqFields)
		if err != nil {
			http.Error(w, "入力に不正があります。", http.StatusInternalServerError)
			log.Printf("Error while json marshaling: %v", err)
			return
		}

		if strings.ContainsAny(reqFields.Title, ".$") {
			http.Error(w, "ルート名にご使用いただけない文字が含まれています。", http.StatusBadRequest)
			return
		}

		//contextに各フィールドの値を追加
		ctx := req.Context()
		ctx = context.WithValue(ctx, "reqFields", reqFields)
		UpdateRoute.ServeHTTP(w, req.WithContext(ctx))
	}
}
