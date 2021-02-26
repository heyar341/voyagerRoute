package routes

import (
	"app/controllers/auth"
	"app/dbhandler"
	"app/model"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type ResponseMsg struct {
	Msg string `json:"msg"`
}

//requestのフィールドを保存する変数
type MultiSearchRequest struct {
	Title  string                 `json:"title" bson:"title"`
	Routes map[string]interface{} `json:"routes" bson:"routes"`
}

func SaveRoutes(w http.ResponseWriter, req *http.Request) {
	//バリデーション完了後のrequestFieldsを取得
	reqFields, ok := req.Context().Value("reqFields").(MultiSearchRequest)
	if !ok {
		http.Error(w, "エラーが発生しました。もう一度操作を行ってください。", http.StatusInternalServerError)
		log.Printf("Error while getting request fields from reuest's context: %v", ok)
		return
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Error(w, "エラーが発生しました。もう一度操作を行ってください。", http.StatusInternalServerError)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}
	userID := user.ID

	//users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。作成時刻はルート名取得時に作成時刻でソートするため
	userDoc := bson.D{{"_id", userID}}
	now := time.Now().UTC()                                             //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	updateField := bson.M{"multi_route_titles." + reqFields.Title: now} //nested fieldsは.(ドット表記)で繋いで書く
	err := dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateField)

	//routes collectionに保存
	routeDocument := bson.D{
		{"user_id", userID},
		{"title", reqFields.Title},
		{"routes", reqFields.Routes},
	}
	_, err = dbhandler.Insert("googroutes", "routes", routeDocument)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msg := ResponseMsg{Msg: "OK"}
	respJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)
}

//ShowAndEditRoutesTplに編集するルートの情報を渡すためのメソッド
func GetRoute(w http.ResponseWriter, title string, userID primitive.ObjectID) string {
	//確認編集したいルートの名前
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

	//DBから取得したルート情報
	type MultiRoute struct {
		ID     primitive.ObjectID     `json:"_id" bson:"_id"`
		UserID primitive.ObjectID     `json:"user_id" bson:"user_id"`
		Title  string                 `json:"title" bson:"title"`
		Routes map[string]interface{} `json:"routes" bson:"routes"`
	}

	var respRoute MultiRoute
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &respRoute)

	//ShowAndEditRoutesTplで使用するためのデータ形式
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

//ルートを更新保存するためのメソッド
func UpdateRoute(w http.ResponseWriter, req *http.Request) {
	//requestのフィールドを保存するstruct
	type RouteUpdateRequest struct {
		ID            primitive.ObjectID     `json:"id" bson:"_id"`
		Title         string                 `json:"title" bson:"title"`
		PreviousTitle string                 `json:"previous_title" bson:"previous_title"`
		Routes        map[string]interface{} `json:"routes" bson:"routes"`
	}
	//バリデーション完了後のrequestFieldsを取得
	reqFields, ok := req.Context().Value("reqFields").(RouteUpdateRequest)
	if !ok {
		http.Error(w, "エラーが発生しました。もう一度操作を行ってください。", http.StatusInternalServerError)
		log.Printf("Error while getting request fields from reuest's context: %v", ok)
		return
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Error(w, "エラーが発生しました。もう一度操作を行ってください。", http.StatusInternalServerError)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}

	userID := user.ID

	//users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。作成時刻はルート名取得時に作瀬時刻でソートするため
	userDoc := bson.M{"_id": userID}
	now := time.Now().UTC() //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	if reqFields.Title == reqFields.PreviousTitle {
		//「タイムスタンプを更新」
		updateField := bson.M{"multi_route_titles." + reqFields.Title: now} //nested fieldsは.(ドット表記)で繋いで書く
		err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateField)

	} else {
		//「元のルート名を削除」
		deleteField := bson.M{"multi_route_titles." + reqFields.PreviousTitle: ""}
		//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
		//公式ドキュメントURL: https://docs.mongodb.com/manual/reference/operator/update/unset/
		err = dbhandler.UpdateOne("googroutes", "users", "$unset", userDoc, deleteField)

		//「新しいルート名とタイムスタンプを追加」
		updateField := bson.M{"multi_route_titles." + reqFields.Title: now} //nested fieldsは.(ドット表記)で繋いで書く
		err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateField)
	}

	//routes collectionに保存
	routeDoc := bson.M{"_id": reqFields.ID}
	updateDoc := bson.D{
		{"title", reqFields.Title},
		{"routes", reqFields.Routes},
	}
	err := dbhandler.UpdateOne("googroutes", "routes", "$set", routeDoc, updateDoc)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msg := ResponseMsg{Msg: "OK"}
	respJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)
}
