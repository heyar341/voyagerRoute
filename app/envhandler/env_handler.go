package envhandler

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnvVal(keyName string) (string, error) {
	//API呼び出しの準備
	err := godotenv.Load(fmt.Sprintf("env/%s.env", os.Getenv("APP_ENV")))
	if err != nil {
		log.Println("Couldn't load env file")
		return "", err
	}
	//envファイルからkeyNameに応じた値を取得
	envVal := os.Getenv(keyName)
	return envVal, nil
}
