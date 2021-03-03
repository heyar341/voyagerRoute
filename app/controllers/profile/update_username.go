package profile

import (
	"app/dbhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"net/url"
)

func UpdateUserName(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		msg := url.QueryEscape("リクエストメソッドが不正です。")
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
	}
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg := url.QueryEscape("エラーが発生しました。もう一度操作を行ってください。")
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}
	userID := user.ID

	newUserName := req.FormValue("username")
	if newUserName == "" {
		msg := url.QueryEscape("ユーザー名は１文字以上入力してください。")
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}

	//user documentを更新
	userDoc := bson.M{"_id": userID}
	updateDoc := bson.D{{"username", newUserName}}
	err := dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateDoc)
	if err != nil {
		msg := url.QueryEscape("エラーが発生しました。もう一度操作を行ってください。")
		http.Redirect(w, req, "/profile/username_edit_form/?msg="+msg, http.StatusInternalServerError)
		log.Printf("Error while saving multi route: %v", err)
		return
	}

	success := url.QueryEscape("ユーザー名を変更しました。")
	http.Redirect(w, req, "/profile/?success="+success, http.StatusSeeOther)
}
