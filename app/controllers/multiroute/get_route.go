package multiroute

import (
	"app/dbhandler"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

//ShowAndEditRoutesTplに編集するルートの情報を渡すための関数
func GetRoute(w http.ResponseWriter, title string, userID primitive.ObjectID) (string, error) {
	//確認編集したいルートの名前
	routeTitle := title

	routeDoc := bson.M{"user_id": userID, "title": routeTitle}
	DBresp, err := dbhandler.Find("googroutes", "routes", routeDoc, nil)
	if err == mongo.ErrNoDocuments {
		log.Printf("Error while getting multi_route document: %v", err)
		return "", err
	}

	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(DBresp)
	if err != nil {
		log.Printf("Error while bson marshaling multi_route document: %v", err)
		return "", err
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

	return string(respJson), nil
}
