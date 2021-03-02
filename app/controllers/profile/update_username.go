package profile

import (
	"app/dbhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func UpdateUserName(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		msg = "リクエストメソッドが不正です。"
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}
	userID := user.ID

	newUserName := req.FormValue("username")

	//user documentを更新
	userDoc := bson.M{"_id": userID}
	updateDoc := bson.D{{"username", newUserName}}
	err := dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateDoc)
	if err != nil {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
	}

	success := "ユーザー名を変更しました。"
	http.Redirect(w, req, "/mypage/?success="+success, http.StatusSeeOther)
}
