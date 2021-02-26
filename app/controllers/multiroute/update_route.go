package multiroute

import (
	"app/model"
	"app/dbhandler"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

//ルートを更新保存するための関数
func UpdateRoute(w http.ResponseWriter, req *http.Request) {
	//バリデーション完了後のrequestFieldsを取得
	reqFields, ok := req.Context().Value("reqFields").(model.RouteUpdateRequest)
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

	var err error
	//編集するルートのユーザー
	userDoc := bson.M{"_id": userID}

	if reqFields.Title == reqFields.PreviousTitle {
		//「タイムスタンプを更新」
		err = updateUsersRouteTitles(userID, reqFields.Title, "$set")
	} else {
		//「元のルート名をuser documentから削除」
		deleteField := bson.M{"multi_route_titles." + reqFields.PreviousTitle: ""}
		//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
		//公式ドキュメントURL: https://docs.mongodb.com/manual/reference/operator/update/unset/
		err = dbhandler.UpdateOne("googroutes", "users", "$unset", userDoc, deleteField)

		//「新しいルート名とタイムスタンプを追加」
		err = updateUsersRouteTitles(userID, reqFields.Title, "$set")
	}

	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
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
	msg := ResponseMsg{Msg: "OK"}
	respJson, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error while json marshaling: %v", err)
	}
	w.Write(respJson)
}
