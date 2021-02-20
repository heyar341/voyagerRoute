package routes

import (
	"app/controllers/auth"
	"app/dbhandler"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseMsg struct {
	Msg string `json:"msg"`
}

type MultiSearchRequest struct {
	Title  string                 `json:"title" bson:"title"`
	Routes map[string]interface{} `json:"routes" bson:"routes"`
}

func SaveRoutes(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}
	//requestのフィールドを保存する変数
	var reqFields MultiSearchRequest
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &reqFields)
	if err != nil {
		http.Error(w, "aa", http.StatusInternalServerError)
	}

	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		c = &http.Cookie{
			Name:  "sessionId",
			Value: "",
		}
	}

	sessionID, err := auth.ParseToken(c.Value)
	if err != nil {
		msg := "セッション情報が不正です。"
		http.Error(w, msg, http.StatusUnauthorized)
		log.Println(err)
		return
	}
	var userID primitive.ObjectID
	if sessionID != "" {
		userID, err = auth.GetLoginUserID(sessionID)
		if err != nil {
			msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
			http.Error(w, msg, http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	document := bson.D{
		{"user_id", userID},
		{"title", reqFields.Title},
		{"routes", reqFields.Routes},
	}
	//DBに保存
	_, err = dbhandler.Insert("googroutes", "routes", document)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msg := ResponseMsg{Msg: "aaa"}
	respJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)
}