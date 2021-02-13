package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html"
	"net/http"
	"net/url"
	"golang.org/x/crypto/bcrypt"
	"app/dbhandler"
)

type userData struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password []byte `json:"password" bson:"password"`
}

func Login(w http.ResponseWriter, req *http.Request){
	if req.Method != "POST"{
		msg := url.QueryEscape("HTTPメソッドが不正です。")
		http.Redirect(w,req,"/register/?msg="+msg,http.StatusSeeOther)
		return
	}
	//ユーザー名をリクエストから取得
	userName := html.EscapeString(req.FormValue("username"))
	if userName == ""{
		msg := url.QueryEscape("ユーザー名を入力してください。")
		http.Redirect(w,req,"/register/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをリクエストから取得
	password := html.EscapeString(req.FormValue("password"))
	if password == ""{
		msg := url.QueryEscape("パスワードを入力してください。")
		http.Redirect(w,req,"/register/?msg="+msg,http.StatusSeeOther)
		return
	}
	//取得するドキュメントの条件
	userDoc := bson.D{{"username",userName}}
	//DBから取得
	resp, err := dbhandler.Find("googroutes", "users", userDoc)
	if err != nil {
		msg := "ユーザー名またはパスワードが正しくありません。"
		http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)
	}
	//DBから取得した値をmarshal
	bsonByte,_ := bson.Marshal(resp)
	var user userData
	//marshalした値をUnmarshalして、userに代入
	bson.Unmarshal(bsonByte, &user)

	storedPass := user.Password
	//DB内のハッシュ化されたパスワードと入力されたパスワードの一致を確認
	err = bcrypt.CompareHashAndPassword(storedPass,[]byte(password))
	//一致しない場合
	if err != nil{
		msg := "ユーザー名またはパスワードが正しくありません。"
		http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)
	}
	//一致した場合
	http.Redirect(w,req,"/",http.StatusSeeOther)

}
