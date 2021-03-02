package profile

import (
	"app/dbhandler"
	"app/model"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func UpdatePassword(w http.ResponseWriter, req *http.Request) {
	//エラーメッセージを定義
	msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
	//Auth middlewareからuserIDを取得
	user, ok := req.Context().Value("user").(model.UserData)
	if !ok {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/email_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while getting userID from reuest's context: %v", ok)
		return
	}

	currPassword := req.FormValue("current-password")
	newPassword := req.FormValue("password")

	//取得するドキュメントの条件
	userDoc := bson.D{{"_id", user.ID}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "users", userDoc, nil)
	if err != nil {
		msg = "パスワードが正しくありません。"
		http.Redirect(w, req, "/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
		return
	}
	//DBから取得した値をmarshal
	bsonByte, err := bson.Marshal(resp)
	if err != nil {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while bson marshaling user document: %v", err)
		return
	}
	var u model.UserData
	//marshalした値をUnmarshalして、userに代入
	err = bson.Unmarshal(bsonByte, &u)
	if err != nil {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while bson unmarshaling :%v", err)
		return
	}

	storedPass := user.Password
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass, []byte(currPassword))
	//一致しない場合
	if err != nil {
		msg = "パスワードが正しくありません。"
		http.Redirect(w, req, "/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
		return
	}
	//一致した場合
	//パスワードをハッシュ化
	securedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 5)
	if err != nil {
		msg = "エラ〜が発生しました。もう一度操作しなおしてください。"
		http.Redirect(w, req, "/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while hashing password :%v", err)
		return
	}

	updateDoc := bson.D{{"password", securedPassword}}
	err = dbhandler.UpdateOne("googroutes", "users", "$set", userDoc, updateDoc)
	if err != nil {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
		log.Printf("Error while hashing password :%v", err)
		return
	}

	if err != nil {
		msg = "エラーが発生しました。もう一度操作を行ってください。"
		http.Redirect(w, req, "/profile/password_edit_form/?msg="+msg, http.StatusSeeOther)
	}

	c, err := req.Cookie("session_id")
	//Cookieを削除
	c.MaxAge = -1
	http.SetCookie(w, c)

	success := "パスワードの変更に成功しました。ログインしてください。"
	http.Redirect(w, req, "/login_form/?success="+success, http.StatusSeeOther)

}
