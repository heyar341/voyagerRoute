package multiroute

import (
	"app/dbhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func DeleteRoute(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		msg = "HTTPメソッドが不正です。"
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		return
	}

	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}
	userID := user.ID

	routeTitle := req.FormValue("title")

	//削除するルートのユーザー
	userDoc := bson.M{"_id": userID}

	//「元のルート名をuser documentから削除」
	deleteField := bson.M{"multi_route_titles." + routeTitle: ""}
	//documentではなく、document内のフィールドを削除する場合、Deleteではなく、Update operatorの$unsetを使って削除する
	//公式ドキュメントURL: https://docs.mongodb.com/manual/reference/operator/update/unset/
	err := dbhandler.UpdateOne("googroutes", "users", "$unset", userDoc, deleteField)
	if err != nil {
		msg = "エラ-が発生しました。もう一度操作をしなおしてください。"
		http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while deleting %v from multi_route_titles: %v", routeTitle, err)
		return
	}

	////routes collectionから削除 エラー解析に使うかもしれないので、削除せずに残しておく
	//routeDoc := bson.D{{"title", routeTitle}, {"user_id", userID}}
	//err = dbhandler.Delete("googroutes", "routes", routeDoc)
	//if err != nil {
	//	msg = "エラ-が発生しました。もう一度操作をしなおしてください。"
	//	http.Redirect(w, req, "/mypage/show_routes/?msg="+msg, http.StatusSeeOther)
	//	log.Printf("Error while deleting route %v: %v",routeTitle, err)
	//	return
	//}

	//レスポンス作成
	success := "ルート「" + routeTitle + "」を削除しました。"
	http.Redirect(w, req, "/mypage/show_routes/?success="+success, http.StatusSeeOther)
	log.Printf("Route [%v] was deleted", routeTitle)
}
