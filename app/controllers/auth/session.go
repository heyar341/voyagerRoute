package auth

import (
	"app/model"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func genNewSession(userID primitive.ObjectID, w http.ResponseWriter) error {
	//固有のセッションIDを作成
	sessionID := uuid.New().String()
	//sessionをDBに保存
	err := model.CreateNewSession(sessionID, userID)
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
		Name:   "session_id",
		Value:  signedStr,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30, //３０日間有効
	}
	http.SetCookie(w, c)

	return nil
}
