package dbhandler

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDB() (*mongo.Client, context.Context, func(), error) {
	path := fmt.Sprintf("env/%s.env", os.Getenv("APP_ENV"))
	env_err := godotenv.Load(path)
	if env_err != nil {
		log.Println("Couldn't load env file")
	}
	//envファイルからDB情報取得
	DB_HEADER := os.Getenv("DB_HEADER")
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	client, err := mongo.NewClient(options.Client().ApplyURI(DB_HEADER + "://" + DB_USER + ":" + DB_PASSWORD + "@" + DB_HOST + DB_PORT))
	if err != nil {
		log.Printf("DB info: %v", err)
		return nil, nil, func() {}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//5秒経っても処理が終了しない場合、強制終了
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Connect DB: %v", err)
		return nil, nil, func() {}, err
	}
	//clientを返してDB操作をできるようにする
	return client, ctx, cancel, nil
}
