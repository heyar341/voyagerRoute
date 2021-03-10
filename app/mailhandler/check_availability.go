package mailhandler

import (
	"app/model"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
)

//メールアドレスが使用可能かのリクエスト
type validEmailRequest struct {
	Email string `json:"email"`
}

func EmailIsAvailable(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" || req.Method != "POST" {
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}
	//requestのフィールドを保存する変数
	var reqFields validEmailRequest
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &reqFields)
	if err != nil {
		http.Error(w, "入力に不正があります。", http.StatusBadRequest)
		log.Printf("Error while json marshaling: %v", err)
		return
	}

	var isValid = false //メールアドレスが使用可能か入れる変数
	_, err = model.FindUser("email", reqFields.Email)
	//ドキュメントがない場合、メールアドレスは使用可能
	if err == mongo.ErrNoDocuments {
		isValid = true
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"valid": isValid})
}
