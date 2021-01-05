package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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
	uName := req.FormValue("username")
	if uName == ""{
		msg := url.QueryEscape("ユーザー名を入力してください。")
		http.Redirect(w,req,"/register/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをリクエストから取得
	password := req.FormValue("password")
	if password == ""{
		msg := url.QueryEscape("パスワードを入力してください。")
		http.Redirect(w,req,"/register/?msg="+msg,http.StatusSeeOther)
		return
	}
	//DBから読み込み
	client, ctx, err := dbhandler.Connect()
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database("googroutes")
	usersCollection := database.Collection("users")

	//DBからのレスポンスを挿入する変数
	var user userData
	err = usersCollection.FindOne(ctx,bson.D{{"username",uName}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Fatal("ドキュメントが見つかりません")
		}
		log.Fatal(err)
	}
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
