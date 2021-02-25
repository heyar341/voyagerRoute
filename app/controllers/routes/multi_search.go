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
	now := time.Now().UTC()                                             //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
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

func GetRoute(w http.ResponseWriter, title string, userID primitive.ObjectID) string {
	routeTitle := title

	routeDoc := bson.M{"user_id": userID, "title": routeTitle}

	DBresp, err := dbhandler.Find("googroutes", "routes", routeDoc, nil)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
	}

	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(DBresp)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
	}

	type MultiRoutes struct {
		ID     primitive.ObjectID     `json:"_id" bson:"_id"`
		UserID primitive.ObjectID     `json:"user_id" bson:"user_id"`
		Title  string                 `json:"title" bson:"title"`
		Routes map[string]interface{} `json:"routes" bson:"routes"`
	}

	var respRoute MultiRoutes
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &respRoute)

	type JSONResp struct {
		ID         primitive.ObjectID     `json:"id" bson:"_id"`
		Title      string                 `json:"title" bson:"title"`
		Routes     map[string]interface{} `json:"routes" bson:"routes"`
		RouteCount int                    `json:"route_count"`
	}
	//レスポンス作成
	fields := JSONResp{ID: respRoute.ID, Title: respRoute.Title, Routes: respRoute.Routes, RouteCount: len(respRoute.Routes)}
	respJson, err := json.Marshal(fields)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}

	return string(respJson)
}

func UpdateRoute(w http.ResponseWriter, req *http.Request) {
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

	//Cookieからセッション情報取得
	_, err = req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		msg := "ログインしてください。"
		http.Error(w, msg, http.StatusUnauthorized)
		log.Printf("Error while getting cookie: %v", err)
		return
	}

	userID, err := auth.GetLoginUserID(req)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while getting loggedin user: %v", err)
		return
	}

	//users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。作成時刻はルート名取得時に作瀬時刻でソートするため
	userDoc := bson.M{"_id": userID}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	if reqFields.Title == reqFields.PreviousTitle {
		updateField := bson.M{"multi_route_titles." + reqFields.Title: now} //nested fieldsは.(ドット表記)で繋いで書く
		err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateField)

	} else {
		//元のルート名を削除
		deleteField := bson.M{"multi_route_titles": reqFields.PreviousTitle}
		//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
		err = dbhandler.UpdateOne("googroutes", "users","$unset", userDoc, deleteField)
		//新しいルート名とタイムスタンプを追加
		updateField := bson.M{"multi_route_titles." + reqFields.Title: now} //nested fieldsは.(ドット表記)で繋いで書く
		err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateField)
	}

	//routes collectionに保存
	routeDoc := bson.M{"_id": reqFields.ID}
	updateDoc := bson.D{
		{"title", reqFields.Title},
		{"routes", reqFields.Routes},
	}
	err = dbhandler.UpdateOne("googroutes", "routes", "$set", routeDoc, updateDoc)
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
