package auth

import (
	"app/dbhandler"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func genNewSession(userID primitive.ObjectID, w http.ResponseWriter) error {
	//固有のセッションIDを作成
	sessionID := uuid.New().String()
	//sessionをDBに保存
	sessionDoc := bson.D{
		{"session_id", sessionID},
		{"user_id", userID},
	}
	_, err := dbhandler.Insert("googroutes", "sessions", sessionDoc)
	if err != nil {
		log.Printf("Error while inserting session data: %v", err)
		return err
	}

	signedStr, err := createToken(sessionID)
	if err != nil {
		log.Printf("Error while creating a tolen: %v", err)
		return err
	}

	//Cookieの設定
	c := &http.Cookie{
		Name:  "session_id",
		Value: signedStr,
		Path:  "/",
		MaxAge: 60*60*24*30,//３０日間有効
	}
	http.SetCookie(w, c)

	return nil
}