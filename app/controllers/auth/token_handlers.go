package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type customClaim struct {
	jwt.StandardClaims
	SessionID string
}

func createToken(sesId string) (string, error) {
	claim := &customClaim{
		StandardClaims: jwt.StandardClaims{
			//30日間有効
			ExpiresAt: time.Now().Add(720 * time.Hour).Unix(),
		},
		SessionID: sesId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil {
		panic("Can't load env file")
	}
	key := []byte(os.Getenv("TOKENIZE_KEY"))
	signedString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Error happend creating a token: %w", err)
	}
	return signedString, nil
}

func ParseToken(sesVal string) (string, error) {
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil {
		panic("Can't load env file")
	}
	key := []byte(os.Getenv("TOKENIZE_KEY"))
	afterVerifToken, err := jwt.ParseWithClaims(sesVal, &customClaim{}, func(token *jwt.Token) (interface{}, error) {
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
