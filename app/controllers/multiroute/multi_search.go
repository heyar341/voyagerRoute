package multiroute

import (
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
	msg := "エラーが発生しました。もう一度操作を行ってください。"
	//バリデーション完了後のrequestFieldsを取得
	reqFields, ok := req.Context().Value("reqFields").(MultiSearchRequest)
	if !ok {
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while getting request fields from reuest's context: %v", ok)
		return
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}
	userID := user.ID

	/*users collectionのmulti_route_titlesフィールドにルート名と作成時刻を追加($set)する。
	作成時刻はルート名取得時に作成時刻でソートするため*/
	err := updateUsersRouteTitles(userID, reqFields.Title, "$set")
	if err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
	}

	//routes collectionに保存
	routeDocument := bson.D{
		{"user_id", userID},
		{"title", reqFields.Title},
		{"routes", reqFields.Routes},
	}
	_, err = dbhandler.Insert("googroutes", "routes", routeDocument)
	if err != nil {
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
	}

	//レスポンス作成
	w.Header().Set("Content-Type", "application/json")
	msgJSON := ResponseMsg{Msg: "OK"}
	respJson, err := json.Marshal(msgJSON)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)
}

//user documentのmult_route_titlesに新しいルート名とタイムスタンプを追加する関数
func updateUsersRouteTitles(userID primitive.ObjectID, routeTitle, operator string) error {
	userDoc := bson.M{"_id": userID}
	now := time.Now().UTC()                                        //MongoDBでは、timeはUTC表記で扱われ、タイムゾーン情報は入れられない
	updateField := bson.M{"multi_route_titles." + routeTitle: now} //nested fieldsは.(ドット表記)で繋いで書く
	err := dbhandler.UpdateOne("googroutes", "users", operator, userDoc, updateField)
	return err
}
