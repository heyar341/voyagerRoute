package routes

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"html/template"
	"io/ioutil"
	"net/http"
	"app/dbhandler"
	"app/controllers/auth"
)

var route_tpl *template.Template


type ResponseMsg struct {
	Msg string
}

type RoutesRequest struct {
	Title string `json:"title" bson:"title"`
	Routes map[string] interface{} `json:"routes" bson:"routes"`
}



func init()  {
	route_tpl = template.Must(template.ParseGlob("templates/route_search/*"))
}
func SaveRoutes(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST"{
		http.Error(w, "HTTPメソッドが不正です。", http.StatusBadRequest)
		return
	}
	//requestのフィールドを保存する変数
	var reqFields RoutesRequest
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body,&reqFields)
	if err != nil {
		http.Error(w, "aa", http.StatusInternalServerError)
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

	sesId,err := auth.ParseToken(c.Value)
	if err != nil {
		msg := "セッション情報が不正です。"
		http.Redirect(w,req,"/?msg="+msg,http.StatusSeeOther)
		log.Println(err)
		return
	}
	var userId primitive.ObjectID
	if sesId != "" {
		userId,err = auth.GetLoginUserID(sesId)
		if err != nil {
			msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
			http.Error(w,msg,http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
	document := bson.D{
		{"user_id",userId},
		{"title",reqFields.Title},
		{"routes",reqFields.Routes},
	}
	//DBに保存
	_, err = dbhandler.Insert("googroutes", "routes", document)
	if err != nil {
		msg := "エラ〜が発生しました。もう一度操作をしなおしてください。"
		http.Error(w,msg,http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	msg := ResponseMsg{Msg: ""}
	respJson ,err := json.Marshal(msg)
	if err != nil{
		http.Error(w,"問題が発生しました。もう一度操作しなおしてください",http.StatusInternalServerError)
	}
	w.Write(respJson)
}
