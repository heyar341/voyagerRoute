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
	"strings"
	"time"
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
		http.Error(w, "入力に不正があります。", http.StatusInternalServerError)
		log.Printf("Error while json marshaling: %v", err)
		return
	}

	if strings.ContainsAny(reqFields.Title, ".$") {
		http.Error(w, "ルート名にご使用いただけない文字が含まれています。", http.StatusBadRequest)
		return
	}

	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		msg := "ログインしてください。"
		http.Error(w, msg, http.StatusUnauthorized)
		log.Printf("Error while getting cookie: %v", err)
		return
	}

	sessionID, err := auth.ParseToken(c.Value)
	if err != nil {
		msg := "セッション情報が不正です。"
		http.Error(w, msg, http.StatusUnauthorized)
		log.Printf("Error while parsing token: %v", err)
		return
	}
	var userID primitive.ObjectID
	if sessionID != "" {
		userID, err = auth.GetLoginUserID(req)
		if err != nil {
			msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
			http.Error(w, msg, http.StatusInternalServerError)
			log.Printf("Error while getting loggedin user: %v", err)
			return
		}
	}

	//users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。作成時刻はルート名取得時に作瀬時刻でソートするため
	userDoc := bson.D{{"_id", userID}}
	now := time.Now().UTC()                                                                //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	updateField := bson.M{"multi_route_titles." + reqFields.Title: now} //nested fieldsは.(ドット表記)で繋いで書く
	err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateField)

	//routes collectionに保存
	document := bson.D{
		{"user_id", userID},
		{"title", reqFields.Title},
		{"routes", reqFields.Routes},
	}
	_, err = dbhandler.Insert("googroutes", "routes", document)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
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
