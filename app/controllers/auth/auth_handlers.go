package auth

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/url"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"app/dbhandler"

)

type userData struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password []byte `json:"password" bson:"password"`
}

func Register(w http.ResponseWriter, req *http.Request){
	if req.Method != "POST"{
		msg := url.QueryEscape("HTTPメソッドが不正です。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//ユーザー名をリクエストから取得
	uName := req.FormValue("username")
	if uName == ""{
		msg := url.QueryEscape("ユーザー名を入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをリクエストから取得
	password := req.FormValue("password")
	if password == ""{
		msg := url.QueryEscape("パスワードを入力してください。")
		http.Redirect(w,req,"/register_form/?msg="+msg,http.StatusSeeOther)
		return
	}
	//パスワードをハッシュ化
	securedPass,err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	//DBに保存
	client, ctx, err := dbhandler.Connect()
	//処理終了後に切断
	defer client.Disconnect(ctx)
	database := client.Database("googroutes")
	usersCollection := database.Collection("users")
	userDocId, err := usersCollection.InsertOne(ctx,bson.D{
		{"username",uName},
		{"password",securedPass},
	})
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	//固有のセッションIDを作成
	sesId := uuid.New().String()
	//DBに保存
	sessionsCollection := database.Collection("sessions")
	_, err = sessionsCollection.InsertOne(ctx,bson.D{
		{"sessionid",sesId},
		{"userid",userDocId},
	})
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	signedStr,err := createToken(sesId)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	//Cookieの設定
	c := &http.Cookie{
		Name: "sessionId",
		Value: signedStr,
	}
	http.SetCookie(w,c)
	http.Redirect(w,req,"/",http.StatusSeeOther)
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

	sesId,err := parseToken(c.Value)
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
		usersCollection := database.Collection("sessions")
		//DBからのレスポンスを挿入する変数
		err = usersCollection.FindOneAndDelete(ctx,bson.D{{"sessionid",sesId}})
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