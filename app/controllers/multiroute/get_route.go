package multiroute

import (
	"app/dbhandler"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

//ShowAndEditRoutesTplに編集するルートの情報を渡すための関数
func getRoute(title string, userID primitive.ObjectID) (string, error) {
	//確認編集したいルートの名前
	routeTitle := title

	routeDoc := bson.M{"user_id": userID, "title": routeTitle}
	r, err := dbhandler.Find("googroutes", "routes", routeDoc, nil)
	if err == mongo.ErrNoDocuments {
		log.Printf("Error while getting multi_route document: %v", err)
		return "", err
	}

	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(r)
	if err != nil {
		log.Printf("Error while bson marshaling multi_route document: %v", err)
		return "", err
	}

	//JSONとして返す型
	type MultiRoute struct {
		ID         primitive.ObjectID     `json:"_id" bson:"_id"`
		Title      string                 `json:"title" bson:"title"`
		Routes     map[string]interface{} `json:"routes" bson:"routes"`
		RouteCount int                    `json:"route_count"` //unmarshal時には値は入らない
	}

	var respRoute MultiRoute
	err = bson.Unmarshal(bsonByte, &respRoute)
	if err != nil {
		log.Printf("Error while bson unmarshaling multi_route document: %v", err)
		return "", err
	}
	//レスポンス作成
	respRoute.RouteCount = len(respRoute.Routes)
	respJson, err := json.Marshal(respRoute)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}

	//JSONのバイナリ形式のままだとtemplateで読み込めないので、stringに変換
	return string(respJson), nil
}
