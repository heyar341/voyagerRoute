package dbhandler

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func connectDB() (*mongo.Client, context.Context, error) {
	env_err := godotenv.Load("env/dev.env")
	if env_err != nil {
		panic("Can't load env file")
	}
	//envファイルからDB情報取得
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + DB_USER + ":" + DB_PASSWORD + "@" + DB_HOST + ":" + DB_PORT))
	if err != nil {
		log.Fatalln("DB info :", err)
		return nil, nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	//7秒経っても処理が終了しない場合、強制終了
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln("Connect DB :", err)
		return nil, nil, err
	}
	//clientを返してDB操作をできるようにする
	return client, ctx, nil
}
