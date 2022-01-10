package auth

import (
	"app/internal/envhandler"
	"app/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type customClaim struct {
	jwt.StandardClaims
	SessionID string
}

func createToken(sessionID string) (string, error) {
	claim := &customClaim{
		StandardClaims: jwt.StandardClaims{
			//30日間有効
			ExpiresAt: time.Now().Add(720 * time.Hour).Unix(),
		},
		SessionID: sessionID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	keyStr, err := envhandler.GetEnvVal("TOKENIZE_KEY")
	if err != nil {
		return "", err
	}
	key := []byte(keyStr)
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Error happend creating a token: %w", err)
	}
	return signedString, nil
}

func ParseToken(sessionValue string) (string, error) {
	keyStr, err := envhandler.GetEnvVal("TOKENIZE_KEY")
	if err != nil {
		return "", err
	}
	key := []byte(keyStr)
	afterVerifToken, err := jwt.ParseWithClaims(sessionValue, &customClaim{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return "", fmt.Errorf("Someone tried hack the site!")
		}
		return key, nil
	})

	if err != nil {
		return "", fmt.Errorf("couldn't parseTokenWithClaim at parseToken: %w", err)
	}
	if !afterVerifToken.Valid {
		return "", fmt.Errorf("セッション情報が不正です。")
	}

	return afterVerifToken.Claims.(*customClaim).SessionID, nil
}

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
