package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"app/dbhandler"
)

func Logout(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Redirect(w,req, "/",http.StatusSeeOther)
	}
	//Cookieからセッション情報取得
	c, err := req.Cookie("sessionId")
	//Cookieが設定されてない場合
	if err != nil {
		c = &http.Cookie{
			Name: "sessionId",
			Value: "",
		}
	}

	sesId,err := ParseToken(c.Value)
	if err != nil {
		msg := "ログインしていません"
		http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)
		log.Println(err)
		return
	}

	if sesId != "" {
		//DBから読み込み
		client, ctx, err := dbhandler.Connect()
		//処理終了後に切断
		defer client.Disconnect(ctx)
		database := client.Database("googroutes")
		sessionsCollection := database.Collection("sessions")
		//DBからのレスポンスを挿入する変数
		var deletedDocument bson.M
		err = sessionsCollection.FindOneAndDelete(ctx,bson.D{{"sessionid",sesId}}).Decode(&deletedDocument)
		if err != nil {
			msg := "エラ〜が発生しました。"
			http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)
			if err == mongo.ErrNoDocuments {
				log.Fatal("Couldn't find a document")
			}
			log.Fatal(err)
			return
		}
	}

	c.MaxAge = -1
	http.SetCookie(w,c)
	msg := "ログアウトしました"
	http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)

}
