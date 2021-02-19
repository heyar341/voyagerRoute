package envhandler

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnvVal(keyName string) string {
	//API呼び出しの準備
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil {
		log.Println("Can't load env file")
	}
	//envファイルからkeyNameに応じた値を取得
	envVal := os.Getenv(keyName)
	return envVal
}
