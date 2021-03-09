package mailhandler

import (
	"app/dbhandler"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
)

func EmailIsAvailable(w http.ResponseWriter, req *http.Request) {
	//メールアドレスが使用可能かのリクエスト
	type ValidEmailRequest struct {
		Email string `json:"email"`
	}

	if req.Method != "POST" {
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}
	//requestのフィールドを保存する変数
	var reqFields ValidEmailRequest
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &reqFields)
	if err != nil {
		http.Error(w, "入力に不正があります。", http.StatusBadRequest)
		log.Printf("Error while json marshaling: %v", err)
		return
	}

	//メールアドレスが使用可能か入れる変数
	var isValid = false
	emailDoc := bson.D{{"email", reqFields.Email}}
	//DBから取得
	_, err = dbhandler.Find("googroutes", "users", emailDoc, nil)
	//ドキュメントがない場合、メールアドレスは使用可能
	if err == mongo.ErrNoDocuments {
		isValid = true
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	type ResponseMsg struct {
		Valid bool `json:"valid"`
	}
	msg := ResponseMsg{Valid: isValid}
	respJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)

}
